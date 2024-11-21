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

type LocalAuth struct{}

var base = "api/user/auth"

func (lo LocalAuth) Join(c *gin.Context) {
	code := c.Query("code")
	types := c.DefaultQuery("type", "user")

	if !lo.WechatIsAction(c) {
		return
	}

	action := "diversion/action"
	if types != "user" {
		action = "action_group"
	}

	if types == "diversion" {
		action = "diversion/action"
	}

	domain := Common.Domain{}.GetTransfer()
	c.HTML(http.StatusOK, "cookie.html", gin.H{
		"code": code,
		//"next_link": fmt.Sprintf("%s/%s/transfer_action/", domain.Domain, base),
		"next_link": fmt.Sprintf("%s/%s/%s/", domain.Domain, base, action),
		"bind_link": fmt.Sprintf("%s/%s/bind_uuid/", Base.AppConfig.HttpHost, base),
		"uuid":      "",
		"type":      types,
		"action":    "join",
	})
}

func (lo LocalAuth) JoinGroup(c *gin.Context) {
	code := c.Query("code")
	types := c.DefaultQuery("type", "user")

	if !lo.WechatIsAction(c) {
		return
	}

	action := "action_group"

	domain := Common.Domain{}.GetTransfer()
	c.HTML(http.StatusOK, "cookie.html", gin.H{
		"code": code,
		//"next_link": fmt.Sprintf("%s/%s/transfer_action/", domain.Domain, base),
		"next_link": fmt.Sprintf("%s/%s/%s/", domain.Domain, base, action),
		"bind_link": fmt.Sprintf("%s/%s/bind_uuid/", Base.AppConfig.HttpHost, base),
		"uuid":      "",
		"type":      types,
		"action":    "join",
	})
}

func (lo LocalAuth) Transfer(c *gin.Context) {
	code := c.Param("code")
	JoinUuid := c.Param("uuid")
	types := c.DefaultQuery("type", "user")

	if !lo.WechatIsAction(c) {
		return
	}

	domain := Common.Domain{}.GetAction()
	c.HTML(http.StatusOK, "cookie.html", gin.H{
		"code":      code,
		"next_link": fmt.Sprintf("%s/%s/action/", domain, base),
		"bind_link": fmt.Sprintf("%s/%s/bind_uuid/", Base.AppConfig.HttpHost, base),
		"type":      types,
		"uuid":      JoinUuid,
		"action":    "transfer",
	})
}

func (lo LocalAuth) Action(c *gin.Context) {
	code := c.Param("code")
	JoinUuid := c.Param("uuid")

	if !lo.WechatIsAction(c) {
		return
	}

	domain := Common.Domain{}.GetAction()
	c.HTML(http.StatusOK, "cookie.html", gin.H{
		"code":      code,
		"next_link": fmt.Sprintf("%s/%s/location?url=%s/%s/show/", domain, base, domain, base),
		"bind_link": fmt.Sprintf("%s/%s/bind_uuid/", Base.AppConfig.HttpHost, base),
		"uuid":      JoinUuid,
		"action":    "action",
	})
}

func (lo LocalAuth) Show(c *gin.Context) {
	code := c.Param("code")
	JoinUuid := c.Param("uuid")

	if !lo.WechatIsAction(c) {
		return
	}

	service, err := Service2.Service{}.Get(code)
	if err != nil {
		c.String(http.StatusOK, "未知客服")
		return
	}

	// 检测账号是否过期
	if service.TimeOut.Unix()-time.Now().Unix() <= 0 {
		Common2.ApiResponse{}.Error(c, "账号已过期", gin.H{})
		return
	}

	// 如果UUID不存在则注册用户并创建cookie、
	userMap := Logic.User{}.CookieUUIDToUser(JoinUuid)

	// 准备绑定的用户 如果cookie里面有uuid 则记录上层UUID的绑定关系 否則創建並注冊
	userModel := Logic.User{}.CheckCookieUUIDToUser(userMap, JoinUuid)

	_ = Service2.ServiceRoom{}.Get(
		userModel, service, c.ClientIP(), Common2.ClientAgentTools{}.GetDrive(c),
		Common2.ClientAgentTools{}.GetDriveInfo(c), "",
	)
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

	link := fmt.Sprintf("%s/api/user/action/detail", action)

	c.HTML(http.StatusOK, "localstorage.html", gin.H{
		"key": "token",
		"val": token,
		"url": link,
	})
}

func (lo LocalAuth) ShowGroup(c *gin.Context) {
	code := c.Param("code")
	JoinUuid := c.Param("uuid")

	if !lo.WechatIsAction(c) {
		return
	}

	group := Common.Group{}.CodeGet(code)
	if group.GroupId == 0 {
		c.String(http.StatusOK, "未知群组")
		return
	}

	service, err := Service2.Service{}.IdGet(group.ServiceId)
	if err != nil {
		c.String(http.StatusOK, "未知客服")
		return
	}

	// 检测账号是否过期
	if service.TimeOut.Unix()-time.Now().Unix() <= 0 {
		Common2.ApiResponse{}.Error(c, "账号已过期", gin.H{})
		return
	}

	// 如果UUID不存在则注册用户并创建cookie、
	userMap := Logic.User{}.CookieUUIDToUser(JoinUuid)

	// 如果cookie里面有uuid 则记录上层UUID的绑定关系 否則創建並注冊
	userModel := Logic.User{}.CheckCookieUUIDToUser(userMap, JoinUuid)

	token := Common2.Tools{}.EncodeToken(userModel.UserId, "user", service.ServiceId, group.GroupId, service.Diversion)
	action := Common.Domain{}.GetAction()
	if action == "" {
		Common2.ApiResponse{}.Error(c, "无法找到action", gin.H{})
		return
	}

	groupUser := Common.Group{}.UserIdGet(group.GroupId, group.ServiceId, userModel.UserId)
	if groupUser.Status == "ban" {
		Common2.ApiResponse{}.Error(c, "[BAN]无法查询数据", gin.H{})
		return
	}

	groupUser.UpdateTime = time.Now()
	go Base.MysqlConn.Model(&Common.Group{}).Updates(&groupUser)

	go Base.MysqlConn.Create(&User2.UserLoginLog{
		UserId:     userModel.UserId,
		ServiceId:  service.ServiceId,
		Ip:         c.ClientIP(),
		Addr:       "",
		CreateTime: time.Now(),
	})
	link := fmt.Sprintf("%s/api/user/action_group/%s#/group", action, "detail")
	c.HTML(http.StatusOK, "localstorage.html", gin.H{
		"key": "token",
		"val": token,
		"url": link,
	})
}

func (LocalAuth) BindUUid(c *gin.Context) {
	newUuid := c.Param("new_uuid")
	currentUuid := c.Param("current_uuid")
	action := c.Query("action")

	// 如果和cookie的用户对不上则绑定一下
	user := Logic.User{}.CookieUUIDToUser(currentUuid)
	if user.UserId == 0 {
		c.String(http.StatusOK, "非法用户.")
		return
	}
	Base.MysqlConn.Create(User2.UserAuthMap{CookieUid: newUuid, UserId: user.UserId, Action: action})
	c.String(http.StatusOK, "ok")

}

func (lo LocalAuth) Location(c *gin.Context) {
	action := c.Query("url")
	c.Redirect(http.StatusTemporaryRedirect, action)
}

func (lo LocalAuth) TransferAction(c *gin.Context) {
	code := c.Param("code")
	JoinUuid := c.Param("uuid")
	domain := Common.Domain{}.GetAction()
	link := fmt.Sprintf("%s/%s/location?url=%s/%s/transfer/%s/%s", domain, base, domain, base, code, JoinUuid)
	link = fmt.Sprintf("%s/%s/action/%s/%s", domain, base, code, JoinUuid)

	c.HTML(http.StatusOK, "but.html", gin.H{"url": link})
}

func (LocalAuth) WechatIsAction(c *gin.Context) bool {
	if Base.AppConfig.Debug {
		return true
	}
	isWeChat := Common2.ClientAgentTools{}.IsWechat(c)
	isDouyin := Common2.ClientAgentTools{}.IsDouYin(c)
	isAilipay := Common2.ClientAgentTools{}.IsAlipay(c)
	isMobile := Common2.ClientAgentTools{}.IsMobile(c)

	if !isWeChat && !isDouyin && !isAilipay {
		c.Redirect(http.StatusTemporaryRedirect, "http://www.baidu.com/")
		return false
	}

	if !isMobile {
		c.Redirect(http.StatusTemporaryRedirect, "http://www.qq.com/")
		return false
	}
	return true
}

func (lo LocalAuth) ActionGroup(c *gin.Context) {
	code := c.Param("code")
	JoinUuid := c.Param("uuid")

	if !lo.WechatIsAction(c) {
		return
	}

	domain := Common.Domain{}.GetAction()

	c.HTML(http.StatusOK, "cookie.html", gin.H{
		"code":      code,
		"next_link": fmt.Sprintf("%s/%s/location?url=%s/%s/show_group/", domain, base, domain, base),
		"bind_link": fmt.Sprintf("%s/%s/bind_uuid/", Base.AppConfig.HttpHost, base),
		"uuid":      JoinUuid,
		"action":    "action",
	})
}
