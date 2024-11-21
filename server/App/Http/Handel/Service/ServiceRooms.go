package Service

import (
	"github.com/gin-gonic/gin"
	"math"
	"server/App/Common"
	"server/App/Http/Request"
	Service3 "server/App/Http/Request/Service"
	"server/App/Http/Response"
	Service4 "server/App/Logic/Service"
	Message2 "server/App/Model/Common"
	Service2 "server/App/Model/Service"
	"server/App/Model/User"
	"server/Base"
	"time"
)

type ServiceRooms struct{}

// @summary 房间-获取用户房间列表
// @tags 房间信息
// @Param token header string true "认证token"
// @Param user_name query string false "用户名"
// @Param type query string false "all 所有 user_no_read  用户未读 server_read 已回复 server_no_read 未回复 top 置顶 black 拉黑"
// @Param page query int true "指定页"
// @Param offset query int true "指定每页数量"
// @Router /api/service/rooms/list [post]
func (ServiceRooms) List(c *gin.Context) {
	var req Service3.ServiceRoomList
	err := c.ShouldBind(&req)

	if err != nil {
		Common.ApiResponse{}.Error(c, "请求繁忙", gin.H{})
		return
	}
	RoleId := Common.Tools{}.GetRoleId(c)
	res := Service4.ServiceRoom{}.List(RoleId, req, true)
	if req.IsClearStatus == 1 {
		Base.WebsocketHub.BindUser(Common.Tools{}.GetWebSocketId(c), 0)
	}

	//序列化
	Common.ApiResponse{}.Success(c, "ok", gin.H{"res": res})
}

func (ServiceRooms) End(c *gin.Context) {
	var req Service3.UserId
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请求繁忙", gin.H{})
		return
	}
	Base.WebsocketHub.BindUser(Common.Tools{}.GetWebSocketId(c), 0)
	Base.MysqlConn.Model(&Service2.ServiceRoom{}).Where("service_id = ? and user_id = ?", Common.Tools{}.GetRoleId(c), req.UserId).Updates(gin.H{"is_delete": 1, "late_msg": ""})
	Base.MysqlConn.Delete(
		&Message2.Message{}, "service_id =  ? and user_id = ?", Common.Tools{}.GetServiceId(c), req.UserId)
	Common.ApiResponse{}.Success(c, "ok", gin.H{"req": req})
}

func (ServiceRooms) Update(c *gin.Context) {
	var req Service3.ServiceRoomDetail
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, err.Error(), gin.H{"req": req})
		return
	}

	serviceId := Common.Tools{}.GetServiceId(c)

	Service4.ServiceRoom{}.Rename(serviceId, req.UserId, req.Rename)

	var findInfo Service2.ServiceRoomDetail
	Base.MysqlConn.Model(&findInfo).Where("service_id = ? and user_id = ?", serviceId, req.UserId).Find(&findInfo)

	updates := gin.H{
		"mobile": req.Mobile,
		"tag":    req.Tag,
	}

	var info Service2.ServiceRoomDetail
	Base.MysqlConn.Model(&info).Where("service_id = ? and user_id = ?", serviceId, req.UserId).Updates(updates)
	Common.ApiResponse{}.Success(c, "ok", gin.H{"info": info})
}

// @summary 房间-获取用户房间详细
// @tags 房间信息
// @Param token header string true "认证token"
// @Param user_id query int true "用户ID"
// @Router /api/service/rooms/detail [post]
func (ServiceRooms) Detail(c *gin.Context) {
	var req Service3.UserId
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请求繁忙", gin.H{})
		return
	}
	serviceId := Common.Tools{}.GetServiceId(c)
	Base.WebsocketHub.BindUser(Common.Tools{}.GetWebSocketId(c), req.UserId)
	//Base.WebsocketHub.BindGroup(Common.Tools{}.GetWebSocketId(c), 0)

	// service未读更新0
	Base.MysqlConn.Model(&Service2.ServiceRoom{}).Where("service_id = ? and user_id = ?", serviceId, req.UserId).Updates(gin.H{"service_no_read": 0})

	// 获取用户信息
	var users Response.UserDetail
	sql := "select service_rooms.`rename`,users.user_id,users.user_name,users.user_head,is_top,service_room_details.drive,service_room_details.drive_info," +
		"service_room_details.ip,service_room_details.map,service_room_details.mobile,service_room_details.tag " +
		"from service_rooms left join users on service_rooms.user_id = users.user_id " +
		"left join service_room_details on service_rooms.user_id = service_room_details.user_id and service_rooms.service_id = service_room_details.service_id " +
		"where service_rooms.service_id = ? and service_rooms.user_id = ?"
	Base.MysqlConn.Raw(sql, serviceId, req.UserId).Scan(&users)

	var UserLoginRecorder []User.UserLoginLog
	Base.MysqlConn.Find(&UserLoginRecorder, "service_id = ? and user_id = ?", serviceId, req.UserId)

	var UserLoginRecorderResp []User.UserLoginLog
	for _, v := range UserLoginRecorder {
		UserLoginRecorderResp = append(UserLoginRecorderResp, User.UserLoginLog{
			Id:         v.Id,
			UserId:     v.UserId,
			ServiceId:  v.ServiceId,
			Ip:         v.Ip,
			Addr:       v.Addr,
			CreateTime: v.CreateTime,
		})
	}
	Common.ApiResponse{}.Success(c, "ok", gin.H{"user": users, "login_log": UserLoginRecorderResp})
}

func (ServiceRooms) Top(c *gin.Context) {
	var req Service3.RoomTop
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请求繁忙", gin.H{})
		return
	}
	serviceId := Common.Tools{}.GetServiceId(c)
	Base.MysqlConn.Model(&Service2.ServiceRoom{}).
		Where("service_id = ? and user_id = ?", serviceId, req.UserId).Update("is_top", req.Top)
	var str string
	if req.Top == 1 {
		str = "置顶成功"
	} else {
		str = "取消置顶"
	}
	Common.ApiResponse{}.Success(c, str, gin.H{"user_id": req.UserId})
}

func (ServiceRooms) Black(c *gin.Context) {
	var req Service3.RoomBlack
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请求繁忙", gin.H{})
		return
	}

	serviceId := Common.Tools{}.GetServiceId(c)
	if req.Type == "ip" {
		Service4.ServiceBlack{}.IpBlack(req.IsBlack, serviceId, req.Ip)
	} else {
		Service4.ServiceBlack{}.UserBlack(req.IsBlack, serviceId, req.Ip, req.UserId)
	}
	Common.ApiResponse{}.Success(c, "ok", gin.H{"user_id": req.UserId})
}

func (ServiceRooms) BlackList(c *gin.Context) {
	var pageReq Request.PageLimit
	err := c.ShouldBind(&pageReq)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请求繁忙", gin.H{})
		return
	}

	ServiceId := Common.Tools{}.GetServiceId(c)

	tel := Base.MysqlConn.Table("service_blacks").
		Select("service_blacks.*,users.user_name,users.user_head").
		Joins("left JOIN users ON service_blacks.user_id = users.user_id").
		Where("service_blacks.service_id = ?", ServiceId)

	// 计算分页和总数
	var allCount int
	tel.Count(&allCount)
	allPage := math.Ceil(float64(allCount) / float64(pageReq.Limit))

	// 获取分页数据
	var list []Response.UserBlackList
	tel.Offset(pageReq.GetOffset()).Limit(pageReq.Limit).Scan(&list)
	res := gin.H{"count": allCount, "page": allPage, "current_page": pageReq.Page, "list": list}
	Common.ApiResponse{}.Success(c, "获取成功", res)
}

func (ServiceRooms) Rename(c *gin.Context) {
	var req Service3.UpdateServiceRoomRename
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	ServiceId := Common.Tools{}.GetServiceId(c)

	Service4.ServiceRoom{}.Rename(ServiceId, req.UserId, req.Rename)
	Common.ApiResponse{}.Success(c, "修改备注成功", gin.H{})
}

func (ServiceRooms) DeleteDay(c *gin.Context) {
	var req Response.DeleteUserDay
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	ServiceId := Common.Tools{}.GetServiceId(c)

	now := time.Now()
	Base.MysqlConn.Model(&Service2.ServiceRoom{}).Where("service_id = ? and update_time >= ?", ServiceId, now.AddDate(0, 0, -req.Day)).Updates(gin.H{"is_delete": 1, "late_msg": ""})
	Base.MysqlConn.Delete(&Message2.Message{}, "service_id = ? and create_time >= ?", ServiceId, now.AddDate(0, 0, -req.Day))

	Common.ApiResponse{}.Success(c, "用户清理成功", gin.H{"req": req, "time": now.AddDate(0, 0, -req.Day)})
}
