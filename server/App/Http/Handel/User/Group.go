package User

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"server/App/Common"
	Common3 "server/App/Http/Request/Common"
	Service3 "server/App/Http/Request/Service"
	Common2 "server/App/Logic/Common"
	Group2 "server/App/Model/Group"
	"server/Base"
)

type Group struct{}

func (Group) Detail(c *gin.Context) {
	roleId := Common.Tools{}.GetServiceId(c)
	group, users := Common2.Group{}.Detail(roleId)
	Common.ApiResponse{}.Success(c, "OK", gin.H{"group": group, "users": users})
}

func (Group) Send(c *gin.Context) {
	var req Common3.GroupSocketMessage
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, req.GetErr(err.(validator.ValidationErrors)), gin.H{})
		return
	}

	serviceId := Common.Tools{}.GetServiceId(c)
	RoleId := Common.Tools{}.GetRoleId(c)
	RoleType := Common.Tools{}.GetRoleType(c)

	// 如果是用户发送消息,检测是否在群组中并且不在黑名单中
	if RoleType == "user" {
		group := Common2.Group{}.GetGroupCache(serviceId)
		groupUser := Common2.Group{}.GetGroupUserCache(serviceId, RoleId)

		if group.Status == "stop" && groupUser.Role == "user" {
			Common.ApiResponse{}.Error(c, "群聊禁止发言", gin.H{})
			return
		}

		if groupUser.ServiceId == 0 {
			Common.ApiResponse{}.Error(c, "你不在群组中", gin.H{})
			return
		}
		if groupUser.IsBlack == 1 {
			Common.ApiResponse{}.Error(c, "你已被拉黑", gin.H{})
			return
		}
		if groupUser.Status == "stop" || groupUser.Status == "ban" {
			Common.ApiResponse{}.Error(c, "您已被禁止发言", gin.H{})
			return
		}

	}

	Common2.Group{}.SendMessage(serviceId, RoleId, RoleType, req)
	Common.ApiResponse{}.Success(c, "ok", gin.H{})
}

func (Group) Update(c *gin.Context) {
	var req Service3.GroupUpdateReq
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, req.GetErr(err.(validator.ValidationErrors)), gin.H{})
		return
	}
	Common2.Group{}.Update(req.GroupId, req)
	Common.ApiResponse{}.Success(c, "OK", gin.H{})
}

func (Group) Message(c *gin.Context) {
	var req Service3.GroupMessageReq
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, req.GetErr(err.(validator.ValidationErrors)), gin.H{})
		return
	}
	serviceId := Common.Tools{}.GetServiceId(c)
	list := Common2.Group{}.MessageList(serviceId, req)
	Common.ApiResponse{}.Success(c, "OK", gin.H{"list": list})
}

func (Group) User(c *gin.Context) {
	var req Service3.UserId
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "user id not find", gin.H{})
		return
	}

	serviceId := Common.Tools{}.GetServiceId(c)
	var groupUser Group2.GroupUser
	Base.MysqlConn.Where("service_id = ? and user_id = ?", serviceId, req.UserId).Find(&groupUser)
	Common.ApiResponse{}.Success(c, "OK", gin.H{"user": groupUser})
}
