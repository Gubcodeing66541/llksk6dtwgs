package User

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"math"
	Common2 "server/App/Common"
	"server/App/Http/Request"
	Service2 "server/App/Http/Request/Service"
	User2 "server/App/Http/Request/User"
	"server/App/Http/Response"
	"server/App/Logic/Common"
	"server/App/Logic/Service"
	Logic "server/App/Logic/User"
	Message2 "server/App/Model/Common"
	Service3 "server/App/Model/Service"
	"server/App/Model/Setting"
	User3 "server/App/Model/User"
	"server/Base"
	"time"
)

type User struct{}

// Token 用户token字符串换取真实token
func (User) Token(c *gin.Context) {
	var req User2.Token
	err := c.ShouldBind(&req)

	if err != nil {
		Common2.ApiResponse{}.Error(c, "请输入需要发送的消息.", gin.H{"token": ""})
		return
	}

	token := Common2.RedisTools{}.GetString(req.Token)
	Common2.ApiResponse{}.Success(c, "解析成功", gin.H{"token": token})
}

// Info 获取客服基本信息
func (User) Info(c *gin.Context) {
	serviceId := Common2.Tools{}.GetServiceId(c)

	// 获取客服信息
	service, err := Service.Service{}.IdGet(serviceId)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "无法获取到客服信息", gin.H{})
		return
	}

	room := Service.ServiceRoom{}.GetByUserId(serviceId, Common2.Tools{}.GetRoleId(c))
	if room.IsBlack == 1 {
		Common2.ApiResponse{}.Ban(c, "无法获取到客服信息", gin.H{})
		return
	}

	// 获取自己的信息
	users := Logic.User{}.UserIdToUser(Common2.Tools{}.GetRoleId(c))

	roleId := Common2.Tools{}.GetRoleId(c)

	// 记录登录日志
	Base.MysqlConn.Create(&User3.UserLoginLog{UserId: roleId, ServiceId: service.ServiceId, Ip: c.ClientIP(), Addr: "", CreateTime: time.Now()})

	// 所有消息已读
	go func() {
		Base.MysqlConn.Model(Message2.Message{}).Where("service_id = ? and user_id = ? and is_read = 0",
			service.ServiceId, users.UserId).Updates(
			gin.H{"is_read": 1})
	}()

	domain := Common.Domain{}.GetServiceBind(service.ServiceId)
	//str := fmt.Sprintf("MEMBER-%s-time:-rand-%d", time.Now(), rand.Intn(99999))
	web := fmt.Sprintf("%s/api/user/ua/join/%s", domain.Domain, service.Code)

	Common2.ApiResponse{}.Success(c, "解析成功", gin.H{
		"users": users,
		"service": Response.ServiceInfo{
			ServiceId:  service.ServiceId,
			Name:       service.Name,
			Head:       service.Head,
			Type:       service.Type,
			Code:       service.Code,
			Diversion:  service.Diversion,
			Web:        web,
			TimeOut:    service.TimeOut.Format("2006-01-02 15:04:05"),
			CreateTime: service.CreateTime,
			UserDetail: service.UserDefault,
		},
		"token": Common2.Tools{}.GetCookieToken(c),
	})
}

// BotMessage 获取智能消息
func (User) BotMessage(c *gin.Context) {
	serviceId := Common2.Tools{}.GetServiceId(c)

	// 获取客服信息
	service, err := Service.Service{}.IdGet(serviceId)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "无法获取到客服信息", gin.H{})
		return
	}

	room := Service.ServiceRoom{}.GetByUserId(serviceId, Common2.Tools{}.GetRoleId(c))
	if room.IsBlack == 1 {
		Common2.ApiResponse{}.Ban(c, "无法获取到客服信息", gin.H{})
		return
	}

	var botMessage []Service3.BotServiceMessage
	Base.MysqlConn.Model(Service3.BotServiceMessage{}).Find(&botMessage, "service_id = ? and status = 'enable'", service.ServiceId)

	Common2.ApiResponse{}.Success(c, "解析成功", gin.H{
		"list": botMessage,
	})
}

// List 消息列表
func (User) List(c *gin.Context) {
	var pageReq Request.PageLimit
	err := c.ShouldBind(&pageReq)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "请提交完整的分页参数", gin.H{})
		return
	}

	roleId := Common2.Tools{}.GetRoleId(c)
	tel := Base.MysqlConn.Model(&Message2.Message{}).Where("service_id = ? and user_id = ? and is_del=0",
		Common2.Tools{}.GetServiceId(c), roleId).Order("id desc")

	var allCount int
	tel.Count(&allCount)

	allPage := float64(0)
	if allCount > 0 {
		allPage = math.Ceil(float64(allCount) / float64(pageReq.GetLimit()))
	}

	var list []Message2.Message
	if pageReq.Id > 0 {
		tel.Where("id < ?", pageReq.Id).Limit(pageReq.GetLimit()).Scan(&list)
	} else {
		tel.Offset(pageReq.GetOffset()).Limit(pageReq.GetLimit()).Scan(&list)
	}

	// 把所有未读变已读
	Base.MysqlConn.Model(&Message2.Message{}).
		Where(" service_id=? and user_id =?", Common2.Tools{}.GetServiceId(c), roleId).
		Update("is_read", 1)

	res := gin.H{"count": allCount, "page": int64(allPage), "current_page": pageReq.Page, "list": list}
	fmt.Println("res", res, "pageReq", pageReq)
	Common2.ApiResponse{}.Success(c, "获取成功", res)
}

// Send 发送消息给客服
func (User) Send(c *gin.Context) {
	var req Service2.UserSendMessage
	err := c.ShouldBind(&req)
	if err != nil || req.Content == "" {
		Common2.ApiResponse{}.Error(c, "请输入需要发送的消息.", gin.H{"erq": req})
		return
	}

	UserId := Common2.Tools{}.GetRoleId(c)
	ServiceId := Common2.Tools{}.GetServiceId(c)
	serviceIsOnline := Base.WebsocketHub.UserIdIsOnline(fmt.Sprintf("%s:%d", "service", ServiceId))
	model := &Message2.Message{
		From: UserId, To: ServiceId, Type: req.Type, Content: req.Content, ServiceId: ServiceId,
		SendRole: "user", CreateTime: time.Now(), IsRead: serviceIsOnline, UserId: UserId, Time: time.Now().Unix()}

	if req.Type != "time" {
		model.IsRead = 1
	}

	Base.MysqlConn.Create(&model)

	if model.Id == 0 {
		Common2.ApiResponse{}.Error(c, "消息发送失败", gin.H{})
		return
	}

	sendMsg := Response.SocketMessage{
		Id:   model.Id,
		From: UserId, To: ServiceId, Type: req.Type, Content: req.Content, ServiceId: ServiceId, Time: time.Now().Unix(),
		SendRole: "user", CreateTime: time.Now().Format("2006-01-02 15:04:05"), IsRead: serviceIsOnline, UserId: UserId,
	}

	if req.Type != "time" {
		LateMsg, LateType := req.Content, req.Type
		update := gin.H{
			"late_role": "user", "late_msg": LateMsg, "update_time": time.Now(), "late_id": model.Id, "is_delete": 0,
			"service_no_read": 0, "user_no_read": 0, "late_type": LateType, "LateUserReadId": model.Id,
		}
		bindUserId := Base.WebsocketHub.GetBindUser(Common2.Tools{}.GetServiceWebSocketId(ServiceId))
		if serviceIsOnline == 0 || bindUserId != UserId {
			update["service_no_read"] = gorm.Expr("service_no_read + ?", 1)
		}
		err = Base.MysqlConn.Table("service_rooms").Where("service_id = ? and user_id = ?", ServiceId, UserId).Updates(update).Error
		fmt.Println("err", err)
		if err != nil {
			fmt.Println("update error", err.Error())
		}
	}

	// 给客服和用户推送
	Common2.ApiResponse{}.SendMsgToService(ServiceId, "message", sendMsg)
	Common2.ApiResponse{}.SendMsgToUser(UserId, "message", sendMsg)

	// 处理离线消息
	Logic.User{}.HandelLeaveMessage(ServiceId, UserId)

	// 更新已读
	go func(ServiceId int, UserId int) {
		Base.MysqlConn.Model(&Message2.Message{}).Where("service_id = ? and user_id = ?", ServiceId, UserId).Updates(gin.H{"is_read": 1})
	}(ServiceId, UserId)

	// ok信息
	Common2.ApiResponse{}.Success(c, "消息发送成功.", gin.H{})
}

// SendBot 客服发送智能消息
func (User) SendBot(c *gin.Context) {
	var req Service2.UserSendMessage
	err := c.ShouldBind(&req)
	if err != nil || req.Content == "" {
		Common2.ApiResponse{}.Error(c, "请输入需要发送的消息.", gin.H{"erq": req})
		return
	}

	// 获取ID
	UserId := Common2.Tools{}.GetRoleId(c)
	ServiceId := Common2.Tools{}.GetServiceId(c)

	// 在线
	serviceIsOnline := Base.WebsocketHub.UserIdIsOnline(fmt.Sprintf("%s:%d", "service", ServiceId))
	model := &Message2.Message{
		From: ServiceId, To: UserId, Type: req.Type, Content: req.Content, ServiceId: ServiceId,
		SendRole: "service", CreateTime: time.Now(), IsRead: serviceIsOnline, UserId: UserId, Time: time.Now().Unix()}

	if req.Type != "time" {
		model.IsRead = 1
	}

	Base.MysqlConn.Create(&model)

	if model.Id == 0 {
		Common2.ApiResponse{}.Error(c, "消息发送失败", gin.H{})
		return
	}

	sendMsg := Response.SocketMessage{
		Id:   model.Id,
		From: ServiceId, To: UserId, Type: req.Type, Content: req.Content, ServiceId: ServiceId, Time: time.Now().Unix(),
		SendRole: "service", CreateTime: time.Now().Format("2006-01-02 15:04:05"), IsRead: serviceIsOnline, UserId: UserId,
	}

	if req.Type != "time" {
		LateMsg, LateType := req.Content, req.Type
		update := gin.H{
			"late_role": "service", "late_msg": LateMsg, "update_time": time.Now(), "late_id": model.Id, "is_delete": 0,
			"service_no_read": 0, "user_no_read": 0, "late_type": LateType, "LateUserReadId": model.Id,
		}
		bindUserId := Base.WebsocketHub.GetBindUser(Common2.Tools{}.GetServiceWebSocketId(ServiceId))
		if serviceIsOnline == 0 || bindUserId != UserId {
			update["service_no_read"] = gorm.Expr("service_no_read + ?", 1)
		}
		err = Base.MysqlConn.Table("service_rooms").Where("service_id = ? and user_id = ?", ServiceId, UserId).Updates(update).Error
		fmt.Println("err", err)
		if err != nil {
			fmt.Println("update error", err.Error())
		}
	}

	// 给客服和用户推送
	Common2.ApiResponse{}.SendMsgToService(ServiceId, "message", sendMsg)
	Common2.ApiResponse{}.SendMsgToUser(UserId, "message", sendMsg)

	// 处理离线消息
	Logic.User{}.HandelLeaveMessage(ServiceId, UserId)

	// 更新已读
	go func(ServiceId int, UserId int) {
		Base.MysqlConn.Model(&Message2.Message{}).Where("service_id = ? and user_id = ?", ServiceId, UserId).Updates(gin.H{"is_read": 1})
	}(ServiceId, UserId)

	// ok信息
	Common2.ApiResponse{}.Success(c, "消息发送成功.", gin.H{})
}

// CreateLive
func (User) CreateLive(c *gin.Context) {
	selviceId := Common2.Tools{}.GetServiceId(c)
	userId := Common2.Tools{}.GetRoleId(c)

	// 生成service的的key和link
	ApiKey := ""
	ApiSecret := ""
	key := Common2.Tools{}.CreateActiveCode(selviceId)
	token := Common2.Tools{}.CreateToken(
		ApiKey, ApiSecret,
		Common2.Tools{}.GetLiveRoomName(int64(selviceId), int64(userId)),
		fmt.Sprintf("S:%d", selviceId),
	)
	link := fmt.Sprintf("%s/api/live/%s", Common.Domain{}.GetLive(), key)
	Common2.RedisTools{}.SetStringByTime(key, token, 60*60)

	// 生成user的的key和link
	userKey := Common2.Tools{}.CreateActiveCode(selviceId)
	userToken := Common2.Tools{}.CreateToken(
		ApiKey, ApiSecret,
		Common2.Tools{}.GetLiveRoomName(int64(selviceId), int64(userId)),
		fmt.Sprintf("U:%d", userId),
	)
	userlink := fmt.Sprintf("%s/api/live/%s", Common.Domain{}.GetLive(), userKey)
	Common2.RedisTools{}.SetStringByTime(userKey, userToken, 60*60)

	// 返回最终结果
	Common2.ApiResponse{}.Success(c, "ok", gin.H{
		"service_key": key, "service_link": link, "user_key": userKey, "user_link": userlink,
	})
}

func (User) Bind(c *gin.Context) {
	var req struct {
		Addr string `json:"addr"`
	}

	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "请求繁忙", gin.H{})
		return
	}

	userId := Common2.Tools{}.GetRoleId(c)
	//Base.MysqlConn.Model(&User2.User{}).Where("user_id = ?", userId).Updates(gin.H{"addr": req.Addr})
	Common2.ApiResponse{}.Success(c, "ok", gin.H{"userId": userId})
}

func (User) IsBind(c *gin.Context) {
	ip := c.ClientIP()
	userId := Common2.Tools{}.GetRoleId(c)
	var userDetail Service3.ServiceRoomDetail
	Base.MysqlConn.Where("user_id = ?", userId).First(&userDetail)

	// 检测IP
	if userDetail.Map == "" {
		var ipModel Common2.Ip
		Base.MysqlConn.Where("ip = ?", ip).First(&ipModel)
		if ipModel.City != "" {
			userDetail.IsBind = 1
			userDetail.Map = ipModel.City
			Base.MysqlConn.Model(&Response.UserDetail{}).Where("user_id = ?", userId).Updates(map[string]interface{}{
				"is_bind": 1, "map": ipModel.City,
			})
		}
	}

	Common2.ApiResponse{}.Success(c, "ok", gin.H{"is_bind": userDetail.IsBind})
}

func (User) CheckoutMessage(c *gin.Context) {
	var req struct {
		MessageId []int64 `json:"message_id"`
	}

	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "请求繁忙", gin.H{})
		return
	}

	// 准备返回的补偿消息
	var msg []Message2.Message

	if len(req.MessageId) == 0 {
		Common2.ApiResponse{}.Success(c, "ok", gin.H{})
		return
	}

	// 查询消息数
	minId := req.MessageId[0]

	for i := 0; i < len(req.MessageId); i++ {
		if minId > req.MessageId[i] {
			minId = req.MessageId[i]
		}
	}

	serviceId := Common2.Tools{}.GetServiceId(c)
	userId := Common2.Tools{}.GetRoleId(c)

	var msgIds []int64
	Base.MysqlConn.Model(&Message2.Message{}).Select("id").
		Where("service_id = ? and user_id = ? and id >= ?", serviceId, userId, minId).Pluck("id", &msgIds)

	msgIdsMap := make(map[int64]bool)
	reqMap := make(map[int64]bool)

	for _, v := range msgIds {
		msgIdsMap[v] = true
	}

	for _, MessageIdV := range req.MessageId {
		reqMap[MessageIdV] = true
	}

	diffId := make([]int64, 0)
	for key, _ := range msgIdsMap {
		if _, ok := reqMap[key]; !ok {
			diffId = append(diffId, key)
		}
	}

	// 补偿消息
	if len(diffId) > 0 {
		Base.MysqlConn.Where("id in (?)", diffId).Find(&msg)
		Common2.ApiResponse{}.Success(c, "ok", gin.H{"messages": msg})
		return
	}

	Common2.ApiResponse{}.Success(c, "ok", gin.H{})
}

func (u User) Loading(c *gin.Context) {
	fmt.Println("1")
	time.Sleep(60 * 10 * time.Second)
}

func (u User) Setting(c *gin.Context) {
	serviceId := Common2.Tools{}.GetServiceId(c)

	var setting Setting.Setting
	Base.MysqlConn.Find(&setting, "service_id=?", serviceId)
	Common2.ApiResponse{}.Success(c, "ok", gin.H{"setting": setting})
}

func (u User) OnlyId(c *gin.Context) {
	var req struct {
		UserId          int    `json:"user_id"`
		ServiceId       int    `json:"service_id"`
		UserToken       string `json:"user_token"`
		FingerprintJSID string `json:"f"`
		CanvasID        string `json:"c"`
	}

	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "请求繁忙", gin.H{})
		return
	}

	Base.MysqlConn.Model(&User3.OnlyId{}).Create(&User3.OnlyId{
		UserId:          req.UserId,
		ServiceId:       req.ServiceId,
		UserToken:       req.UserToken,
		FingerprintJSID: req.FingerprintJSID,
		CanvasID:        req.CanvasID,
		CreateTime:      time.Now(),
		IP:              c.ClientIP(),
	})
	Common2.ApiResponse{}.Success(c, "ok", gin.H{})

}
