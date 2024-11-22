package Service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"server/App/Common"
	Service3 "server/App/Http/Request/Service"
	"server/App/Http/Response"
	Service2 "server/App/Logic/Service"
	Message2 "server/App/Model/Common"
	Service4 "server/App/Model/Service"
	"server/Base"
	"time"
)

type Message struct{}

func (Message) SendToUser(c *gin.Context) {
	var req Service3.ServiceSendMessage
	err := c.ShouldBind(&req)
	if err != nil || req.UserId == 0 || req.Content == "" {
		Common.ApiResponse{}.Error(c, "请输入需要发送的消息.", gin.H{"erq": req})
		return
	}

	RoleId := Common.Tools{}.GetRoleId(c)
	err = Service2.Message{}.SendToUser(RoleId, req.UserId, req.Type, req.Content, true)
	if err != nil {
		Common.ApiResponse{}.Error(c, err.Error(), gin.H{"erq": req})
	}

	// 返回OK信息
	Common.ApiResponse{}.Success(c, "消息发送成功.", gin.H{})

}

func (Message) SendAll(c *gin.Context) {
	var req Service3.ServiceSendMessageGroup
	userList := c.PostFormArray("user_id")

	err := c.ShouldBind(&req)
	if err != nil || req.Content == "" || req.Type == "" {
		Common.ApiResponse{}.Error(c, "请输入需要发送的消息.", gin.H{"erq": req, "userList": userList})
		return
	}

	// 开启协程循环发送
	RoleId := Common.Tools{}.GetRoleId(c)
	go func(RoleId int, req Service3.ServiceSendMessageGroup) {
		// 查询非黑名单用户列表
		var i = 0
		for _, Item := range req.UserId {
			i++
			err = Service2.Message{}.SendToUser(RoleId, Item, req.Type, req.Content, true)
			if err != nil {
				fmt.Println("send Error", err.Error())
			}
			if i >= 50 {
				i = 0
				time.Sleep(50 * time.Millisecond)
			}
		}
	}(RoleId, req)

	Common.ApiResponse{}.Success(c, "群发成功", gin.H{"erq": req})

	return
}

func (Message) List(c *gin.Context) {
	var req Service3.MsgList
	err := c.ShouldBind(&req)
	if err != nil || req.UserId == 0 {
		Common.ApiResponse{}.Error(c, "用户不存在", gin.H{})
		return
	}

	roleId := Common.Tools{}.GetRoleId(c)
	var message []Response.Message
	tel := Base.MysqlConn.Model(&Message2.Message{}).Where("service_id = ? and user_id = ? and is_del = 0 and id < ?", roleId, req.UserId, req.Id)

	if req.Id > 0 {
		tel = tel.Where("id < ?", req.Id) // 查找小于指定 id 的记录
	}

	// 计算分页和总数
	var allCount int
	tel.Count(&allCount)
	allPage := math.Ceil(float64(allCount) / float64(req.Offset))

	// 获取分页数据
	tel = Base.MysqlConn.Raw("select * from (select * from messages where service_id = ? and user_id = ?  and is_del = 0 order by id desc limit ? offset ? )  t order by id asc",
		roleId, req.UserId, req.Offset, (req.Page-1)*req.Offset)
	tel.Scan(&message)
	Common.ApiResponse{}.Success(c, "OK", gin.H{"count": allCount, "page": allPage, "current_page": req.Page, "list": message})
}

func (Message) Update(c *gin.Context) {
	var req Service3.UpdateServiceDetail
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "用户不存在", gin.H{})
		return
	}

	roleId := Common.Tools{}.GetRoleId(c)
	Base.MysqlConn.Model(&Service4.Service{}).Where("service_id = ?", roleId).Updates(req)
	Common.ApiResponse{}.Success(c, "ok", gin.H{})
}

func (Message) RemoveMessage(c *gin.Context) {
	var req Service3.RemoveMsg
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "填写完整信息", gin.H{})
		return
	}

	//逻辑删除消息
	//Base.MysqlConn.Unscoped().Delete(&Message2.Message{}, "id = ?", req.Id)
	Base.MysqlConn.Model(&Message2.Message{}).Where("id = ?", req.Id).Updates(
		&Message2.Message{IsDel: 1})

	//是不是最近消息

	var lateMessage Service4.ServiceRoom
	Base.MysqlConn.Find(&lateMessage, "late_id = ? ", req.Id)
	RoleId := Common.Tools{}.GetServiceId(c)
	if req.Id == lateMessage.LateId {
		Base.MysqlConn.Model(&Service4.ServiceRoom{}).Where("service_id = ? and user_id = ?", RoleId, req.UserId).Updates(gin.H{"late_msg": "你撤回了一条消息", "late_type": "text"})
	}

	Common.ApiResponse{}.SendMsgToService(RoleId, "remove", req)
	Common.ApiResponse{}.SendMsgToUser(req.UserId, "remove", req)

	Common.ApiResponse{}.Success(c, "ok", gin.H{"req": req})
}

func (Message) ClearMessage(c *gin.Context) {

	var req Service3.RemoveMsg
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "填写完整信息", gin.H{})
		return
	}
	RoleId := Common.Tools{}.GetServiceId(c)
	//Base.MysqlConn.Delete(&Message2.Message{}, "service_id = ? and user_id = ?", RoleId, req.UserId)

	Base.MysqlConn.Model(&Message2.Message{}).Where("service_id = ? and user_id = ?", RoleId, req.UserId).Updates(
		&Message2.Message{IsDel: 1})

	Base.MysqlConn.Model(&Service2.ServiceRoom{}).Where("service_id = ? and user_id = ?", RoleId, req.UserId).Updates(gin.H{
		"LateId":         0,
		"LateType":       "",
		"LateMsg":        "",
		"LateRole":       "",
		"LateUserReadId": 0,
		"UserNoRead":     0,
		"ServiceNoRead":  0,
	})
	Common.ApiResponse{}.SendMsgToService(RoleId, "clear", req)
	Common.ApiResponse{}.SendMsgToUser(req.UserId, "clear", req)
	Common.ApiResponse{}.Success(c, "ok", gin.H{"req": req})
}

func (Message) ClearLateMessage(c *gin.Context) {
	var req Service3.RemoveLateMsg
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "填写完整信息", gin.H{})
		return
	}

	where := " service_id = ?"
	if req.Type == "week" {
		// 7 * 86400 = 7天
		where += fmt.Sprintf(" and created_at < %d", time.Now().Unix()-7*86400)
	}
	if req.Type == "month" {
		// 30 * 86400 = 30天
		where += fmt.Sprintf(" and created_at < %d", time.Now().Unix()-30*86400)
	}

	RoleId := Common.Tools{}.GetServiceId(c)
	Base.MysqlConn.Delete(&Message2.Message{}, where, RoleId)
	Base.MysqlConn.Model(&Service4.ServiceRoom{}).Where(where, RoleId).Updates(gin.H{"is_delete": 1})
	Common.ApiResponse{}.Success(c, "ok", gin.H{"req": req})
}
