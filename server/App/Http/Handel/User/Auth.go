package User

import (
	"github.com/gin-gonic/gin"
	"net/http"
	Common2 "server/App/Common"
	Service3 "server/App/Http/Request/Service"
	Message2 "server/App/Model/Common"
	"server/App/Model/Service"
	"server/Base"
	"sync"
	"time"
)

type Auth struct{}

// 落地
func (Auth) Action(c *gin.Context) {
	// 通过如果token的code存在则检测
	tokenName := c.Param("token")

	// token 不存在用CookieToken
	if tokenName == "" {
		Common2.ApiResponse{}.Error(c, "未授权", gin.H{})
		return
	}

	token := Common2.RedisTools{}.GetString(tokenName)
	//Common2.RedisTools{}.SetString(tokenName, "")

	var userAuthToken Common2.UserAuthToken
	err := Common2.Tools{}.DecodeToken(token, &userAuthToken)
	if err != nil {
		c.HTML(http.StatusOK, "pay.html", gin.H{"token": token})
		return
	}

	var black Service.ServiceBlack
	Base.MysqlConn.Find(
		&black,
		"service_id = ?  and user_id = ?",
		userAuthToken.ServiceId, userAuthToken.RoleId)
	if black.Id != 0 {
		c.String(200, "..")
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done() // 确保 Goroutine 完成后减少计数
		Base.MysqlConn.Model(&Service3.ServiceRoomDetail{}).Where("service_id = ? and user_id = ?",
			userAuthToken.ServiceId, userAuthToken.RoleId).Updates(
			gin.H{"ip": c.ClientIP()})

		// 上线更新update时间和未读
		update := gin.H{"update_time": time.Now(), "user_no_read": 0, "late_user_read_id": 0, "is_delete": 0, "late_ip": c.ClientIP()}
		Base.MysqlConn.Model(&Service.ServiceRoom{}).
			Where("user_id = ? and service_id = ? ", userAuthToken.RoleId, userAuthToken.ServiceId).
			Updates(update)
	}()

	// 所有消息已读
	Base.MysqlConn.Model(Message2.Message{}).Where("service_id = ? and user_id = ? and is_read = 0",
		userAuthToken.ServiceId, userAuthToken.RoleId).Updates(
		gin.H{"is_read": 1})

	c.HTML(http.StatusOK, "index.html", gin.H{"token": token})

	wg.Wait()
	return
}

// ActionGroup GROUP
func (Auth) ActionGroup(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
	return
}
