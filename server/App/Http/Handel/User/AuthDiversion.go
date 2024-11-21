package User

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	Common2 "server/App/Common"
	"server/App/Logic/Common"
	Service2 "server/App/Logic/Service"
	"server/App/Logic/User"
	Message2 "server/App/Model/Common"
	"server/App/Model/Service"
	User2 "server/App/Model/User"
	"server/Base"
	"time"
)

type AuthDiversion struct{}

func (AuthDiversion) Action(c *gin.Context) {
	code := c.Param("code")
	JoinUuid := c.Param("uuid")

	status := LocalAuth{}.WechatIsAction(c)
	if !status {
		return
	}

	domain := Common.Domain{}.GetAction()
	c.HTML(http.StatusOK, "cookie.html", gin.H{
		"code":      code,
		"next_link": fmt.Sprintf("%s/%s/location?url=%s/%s/diversion/show/", domain, base, domain, base),
		"bind_link": fmt.Sprintf("%s/%s/bind_uuid/", Base.AppConfig.HttpHost, base),
		"uuid":      JoinUuid,
		"action":    "action",
		"key":       Base.AppConfig.IpRegistryKey,
	})
}

func (AuthDiversion) Show(c *gin.Context) {
	code := c.Param("code")
	JoinUuid := c.Param("uuid")

	// 如果UUID不存在则注册用户并创建cookie、
	userMap := Logic.User{}.CookieUUIDToUser(JoinUuid)

	// 准备绑定的用户 如果cookie里面有uuid 则记录上层UUID的绑定关系 否則創建並注冊
	userModel := Logic.User{}.CheckCookieUUIDToUser(userMap, JoinUuid)

	status := LocalAuth{}.WechatIsAction(c)
	if !status {
		return
	}
	service, err := Service2.Service{}.GetByDiversion(code, userModel.UserId)
	if err != nil {
		c.String(http.StatusOK, "未知客服")
		return
	}

	// 检测账号是否过期
	if service.TimeOut.Unix()-time.Now().Unix() <= 0 {
		Common2.ApiResponse{}.Error(c, "账号已过期", gin.H{})
		return
	}

	_ = Service2.ServiceRoom{}.Get(
		userModel, service, c.ClientIP(), Common2.ClientAgentTools{}.GetDrive(c),
		Common2.ClientAgentTools{}.GetDriveInfo(c), "")
	token := Common2.Tools{}.EncodeToken(userModel.UserId, "user", service.ServiceId, 0, service.Diversion)
	action := Common.Domain{}.GetAction()
	if action == "" {
		Common2.ApiResponse{}.Error(c, "无法找到action", gin.H{})
		return
	}

	ip := c.ClientIP()
	var black Service.ServiceBlack
	Base.MysqlConn.Find(
		&black,
		"(service_id = ? and type='ip' and ip = ?) or (service_id = ? and type='user' and user_id = ?)",
		service.ServiceId, ip, service.ServiceId, userModel.UserId)
	if black.Id != 0 {
		c.String(200, "..")
		return
	}

	// 上线更新update时间和未读
	update := gin.H{"update_time": time.Now(), "user_no_read": 0, "late_user_read_id": 0, "is_delete": 0, "late_ip": c.ClientIP()}
	Base.MysqlConn.Model(&Service.ServiceRoom{}).
		Where("user_id = ? and service_id = ? ", userModel.UserId, service.ServiceId).
		Updates(update)

	go Base.MysqlConn.Model(&Service.ServiceRoomDetail{}).Where("service_id = ? and user_id = ?",
		service.ServiceId, userModel.UserId).Updates(
		gin.H{"ip": c.ClientIP()})

	go Base.MysqlConn.Create(&User2.UserLoginLog{
		UserId:     userModel.UserId,
		ServiceId:  service.ServiceId,
		Ip:         c.ClientIP(),
		Addr:       "",
		CreateTime: time.Now(),
	})

	// 所有消息已读
	Base.MysqlConn.Model(Message2.Message{}).Where("service_id = ? and user_id = ? and is_read = 0",
		service.ServiceId, userModel.UserId).Updates(
		gin.H{"is_read": 1})

	link := fmt.Sprintf("%s/api/user/auth/diversion/html", action)

	c.HTML(http.StatusOK, "localstorage.html", gin.H{
		"key": "token",
		"val": token,
		"url": link,
	})
}

func (AuthDiversion) Html(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
	return
}
