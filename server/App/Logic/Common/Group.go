package Common

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	Common2 "server/App/Common"
	"server/App/Http/Request/Common"
	Service3 "server/App/Http/Request/Service"
	GroupModel "server/App/Model/Group"
	"server/App/Model/User"
	"server/Base"
	"strings"
	"time"
)

type Group struct{}

func (Group) SendMessage(serviceId int, RoleId int, RoleType string, req Common.GroupSocketMessage) {
	req.Content = strings.TrimSpace(req.Content)
	groupMessage := GroupModel.GroupMessage{
		ServiceId:  serviceId,
		GroupId:    req.GroupId,
		SendRoleId: RoleId,
		SendRole:   RoleType,
		Type:       req.Type,
		Content:    req.Content,
		ReadCnt:    0,
		CreateTime: time.Now(),
		Time:       time.Now().Unix(),
	}

	// 如果不是仅推送消息，那么就要保存到数据库
	if req.IsJustPush == 0 {
		Base.MysqlConn.Create(&groupMessage)
	}

	ServiceIsOnline := Base.WebsocketHub.UserIdIsOnline(Common2.Tools{}.GetServiceWebSocketId(serviceId))
	if ServiceIsOnline == 1 {
		Base.MysqlConn.Model(&GroupModel.Group{}).Where("service_id = ?", serviceId).Updates(gin.H{
			"late_msg":    req.Content,
			"no_read_cnt": 0,
			"late_time":   time.Now(),
			"update_time": time.Now(),
		})
	} else {
		Base.MysqlConn.Model(&GroupModel.Group{}).Where("service_id = ?", serviceId).Updates(gin.H{
			"late_msg":    req.Content,
			"no_read_cnt": gorm.Expr("no_read_cnt + ?", 1),
			"late_time":   time.Now(),
			"update_time": time.Now(),
		})
	}

	Common2.ApiResponse{}.SendMsgToGroup(groupMessage)
}

func (Group) Detail(roleId int) (GroupModel.Group, []GroupModel.GroupUser) {
	var group GroupModel.Group
	Base.MysqlConn.Find(&group, "service_id = ?", roleId)

	var users []GroupModel.GroupUser
	Base.MysqlConn.Find(&users, "service_id = ?", roleId)
	return group, users
}

func (Group) Update(groupId int, req Service3.GroupUpdateReq) {
	Base.MysqlConn.Model(&GroupModel.Group{}).Where("group_id = ?", groupId).Updates(req)
}

func (Group) CodeGet(code string) GroupModel.Group {
	var group GroupModel.Group
	Base.MysqlConn.Find(&group, "code = ?", code)
	return group
}

// UserIdGet 这里需要做判断，如果就没有加入群组则自动加入
func (Group) UserIdGet(groupId int, serviceId int, userId int) GroupModel.GroupUser {
	var group GroupModel.GroupUser
	Base.MysqlConn.Find(&group, "service_id = ? and user_id = ?", serviceId, userId)

	if group.Id == 0 {
		user := User.User{}
		Base.MysqlConn.Find(&user, "user_id = ?", userId)
		group = GroupModel.GroupUser{
			GroupId:    groupId,
			UserId:     userId,
			Rename:     "",
			Name:       user.UserName,
			Head:       user.UserHead,
			Role:       "user",
			Status:     "none",
			IsBlack:    0,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
			ServiceId:  serviceId,
		}
		Base.MysqlConn.Create(&group)
	}
	return group
}

// MessageList 消息
func (Group) MessageList(roleId int, req Service3.GroupMessageReq) []GroupModel.GroupMessage {
	var group []GroupModel.GroupMessage
	Base.MysqlConn.Where("service_id = ?", roleId).Order("id desc").Limit(req.GetLimit()).Offset(req.GetOffset()).Find(&group)
	return group
}

// DelMessage 删除群聊
func (Group) DelMessage(id string) {
	Base.MysqlConn.Delete(&GroupModel.GroupMessage{}, "id = ?", id)
}

// 判断用户是否被禁言
func (Group) GetGroupUserCache(ServiceId int, userId int) GroupModel.GroupUser {
	var GroupUser GroupModel.GroupUser
	groupUserKey := fmt.Sprintf("group_user-serviceId:%d-user_Id:%d-status", ServiceId, userId)

	userStr := Common2.RedisTools{}.GetString(groupUserKey)
	if userStr == "" {
		Base.MysqlConn.Find(&GroupUser, "service_id = ? and user_id = ?", ServiceId, userId)
		userStr = GroupUser.Status
		userStrByte, _ := json.Marshal(GroupUser)
		Common2.RedisTools{}.SetString(groupUserKey, string(userStrByte))
	} else {
		_ = json.Unmarshal([]byte(userStr), &GroupUser)
	}

	return GroupUser
}

func (Group) ClearGroupUseCache(ServiceId int, userId int) {
	groupUserKey := fmt.Sprintf("group_user-serviceId:%d-user_Id:%d-status", ServiceId, userId)
	Common2.RedisTools{}.Det(groupUserKey)
}

// 判断用户是否被禁言
func (Group) GetGroupCache(ServiceId int) GroupModel.Group {
	var group GroupModel.Group
	groupUserKey := fmt.Sprintf("group-serviceId:%d", ServiceId)
	userStr := Common2.RedisTools{}.GetString(groupUserKey)
	if userStr == "" {
		Base.MysqlConn.Find(&group, "service_id = ?", ServiceId)
		userStr = group.Status
		userStrByte, err := json.Marshal(group)
		if err != nil {
			fmt.Println("err !----------------------", err.Error())
		}
		Common2.RedisTools{}.SetString(groupUserKey, string(userStrByte))
	} else {
		_ = json.Unmarshal([]byte(userStr), &group)
	}

	return group
}

func (Group) ClearGroupCache(ServiceId int) {
	groupUserKey := fmt.Sprintf("group-serviceId:%d", ServiceId)
	Common2.RedisTools{}.Det(groupUserKey)
}
