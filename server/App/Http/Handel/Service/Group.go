package Service

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"server/App/Common"
	Common3 "server/App/Http/Request/Common"
	Service3 "server/App/Http/Request/Service"
	Common2 "server/App/Logic/Common"
	Service2 "server/App/Logic/Service"
	Group2 "server/App/Model/Group"
	"server/Base"
	"strconv"
)

type Group struct{}

func (Group) UpdateUser(c *gin.Context) {
	var req Service3.UpdateGroupUserReq
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, req.GetErr(err.(validator.ValidationErrors)), gin.H{})
		return
	}
	roleId := Common.Tools{}.GetRoleId(c)
	err = Service2.Group{}.UpdateUser(roleId, req)
	if err != nil {
		Common.ApiResponse{}.Error(c, err.Error(), gin.H{})
		return
	}

	// 清理cache
	Common2.Group{}.ClearGroupUseCache(roleId, req.UserId)
	Common.ApiResponse{}.Success(c, "OK", gin.H{})
}

func (Group) Detail(c *gin.Context) {
	roleId := Common.Tools{}.GetRoleId(c)
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

	Common2.Group{}.SendMessage(serviceId, RoleId, RoleType, req)
	Common.ApiResponse{}.Success(c, "ok", gin.H{})
}

func (Group) User(c *gin.Context) {

	var req Service3.UserId
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "user id not find", gin.H{})
		return
	}

	roleId := Common.Tools{}.GetRoleId(c)
	var groupUser Group2.GroupUser
	Base.MysqlConn.Where("service_id = ? and user_id = ?", roleId, req.UserId).Find(&groupUser)
	Common.ApiResponse{}.Success(c, "OK", gin.H{"user": groupUser})
}

func (Group) Update(c *gin.Context) {
	var req Service3.UpdateGroupReq
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, req.GetErr(err.(validator.ValidationErrors)), gin.H{})
		return
	}
	roleId := Common.Tools{}.GetRoleId(c)
	Base.MysqlConn.Model(&Group2.Group{}).Where("service_id = ?", roleId).Updates(req)

	Common2.Group{}.ClearGroupCache(roleId)
	Common.ApiResponse{}.Success(c, "OK", gin.H{})
}

func (Group) UserList(c *gin.Context) {
	var req Service3.GroupUserSearchReq
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, req.GetErr(err.(validator.ValidationErrors)), gin.H{})
		return
	}

	roleId := Common.Tools{}.GetRoleId(c)
	var groupUser []Group2.GroupUser

	if req.UserName == "" {
		Base.MysqlConn.Where("service_id = ?", roleId).Find(&groupUser)
	} else {
		searchName := "%" + req.UserName + "%"
		Base.MysqlConn.Where("service_id = ? and (name like ? or `rename` like ?)", roleId, searchName, searchName).Find(&groupUser)
	}

	Common.ApiResponse{}.Success(c, "OK", gin.H{"users": groupUser})
}

func (Group) Remove(c *gin.Context) {

	var req Service3.GroupDeleteMsgReq
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, req.GetErr(err.(validator.ValidationErrors)), gin.H{})
		return
	}

	roleId := Common.Tools{}.GetRoleId(c)
	roleType := Common.Tools{}.GetRoleType(c)
	Base.MysqlConn.Delete(&Group2.GroupMessage{}, "service_id = ? and id =  ?", roleId, req.Id)

	group, _ := Common2.Group{}.Detail(roleId)

	//roleid 转string
	stringRoleId := strconv.Itoa(req.Id)

	// 群聊推送
	Common2.Group{}.SendMessage(roleId, roleId, roleType, Common3.GroupSocketMessage{
		Type:       "remove",
		Content:    stringRoleId,
		GroupId:    group.GroupId,
		IsJustPush: 1,
	})

	Common.ApiResponse{}.Success(c, "ok", gin.H{})
}
