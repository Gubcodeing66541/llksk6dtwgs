package Common

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	GroupModel "server/App/Model/Group"
	"server/Base"
)

type ApiResponse struct{}

// Success 正确消息
func (ApiResponse) Success(c *gin.Context, msg string, data gin.H) bool {
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": msg, "data": data})
	return true
}

// 错误消息
func (ApiResponse) Error(c *gin.Context, msg string, data gin.H) bool {
	c.JSON(http.StatusOK, gin.H{"code": 300, "msg": msg, "data": data})
	return true
}

// 错误消息
func (ApiResponse) Ban(c *gin.Context, msg string, data gin.H) bool {
	c.JSON(http.StatusOK, gin.H{"code": 302, "msg": msg, "data": data})
	return true
}

// NoAuth 无权限操作
func (ApiResponse) NoAuth(c *gin.Context, data gin.H) bool {
	c.JSON(http.StatusOK, gin.H{"code": 301, "msg": "无操作权限.", "data": data})
	return true
}

// SendMsgToUser 给用户推送socket信息
func (ApiResponse) SendMsgToUser(userId int, typeStr string, content interface{}) {
	userIdName := fmt.Sprintf("%s:%d", "user", userId)
	res := gin.H{"type": typeStr, "content": content}
	resStr, err := json.Marshal(res)
	if err != nil {
		return
	}
	Base.WebsocketHub.SendToUserId(userIdName, resStr)
	return
}

// SendMsgToService 给客服推送socket信息
func (ApiResponse) SendMsgToService(serverId int, typeStr string, content interface{}) {
	userIdName := fmt.Sprintf("%s:%d", "service", serverId)
	res := gin.H{"type": typeStr, "content": content}
	resStr, err := json.Marshal(res)
	if err != nil {
		return
	}
	Base.WebsocketHub.SendToUserId(userIdName, resStr)
}

// SendMsgToGroup 给客服推送socket信息
func (ApiResponse) SendMsgToGroup(msg GroupModel.GroupMessage) {
	res, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("SendMsgToGroupError", err.Error())
		return
	}
	err = NsqHelp{}.Push("group", res)
	if err != nil {
		fmt.Println("NsqHelpPushError", err.Error())
		return
	}
}
