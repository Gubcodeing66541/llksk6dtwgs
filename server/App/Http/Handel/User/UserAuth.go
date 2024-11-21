package User

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server/App/Common"
	Common2 "server/App/Logic/Common"
	Service2 "server/App/Logic/Service"
	Logic "server/App/Logic/User"
	"server/Base"
	"time"
)

type UserAuth struct{}

func (UserAuth) Join(c *gin.Context) {
	// 获取激活码
	code := c.Param("code")

	// 如果没有激活码，直接跳转到百度
	if code == "" {
		c.Redirect(302, "https://www.baidu.com")
		return
	}

	// 获取落地
	action := Common2.Domain{}.GetAction()

	// 生成跳转链接
	link := fmt.Sprintf("%s/api/user/ua/to/%s", action, code)

	// 直接临时跳转
	c.Redirect(302, link)
}

func (UserAuth) To(c *gin.Context) {
	// 获取激活码
	code := c.Param("code")

	// 如果没有激活码，直接跳转到百度
	if code == "" {
		c.Redirect(302, "https://www.baidu.com")
		return
	}

	c.HTML(200, "to.html", gin.H{"code": code, "key": Base.AppConfig.IpRegistryKey})
}

func (UserAuth) Register(c *gin.Context) {
	var req struct {
		Code string `json:"code" binding:"required"`
		City string `json:"city"`
	}
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "未知", gin.H{})
		return
	}

	// 客服
	service, e := Service2.Service{}.Get(req.Code)
	if e != nil {
		Common.ApiResponse{}.Error(c, "未知客服", gin.H{})
		return
	}

	// 检测账号是否过期
	if service.TimeOut.Unix()-time.Now().Unix() <= 0 {
		Common.ApiResponse{}.Error(c, "账号已过期", gin.H{})
		return
	}

	// 用户注册
	username := fmt.Sprintf("%s", Common.Tools{}.GetRename())
	userModel := Logic.User{}.CreateUser("", username, Common.Tools{}.GetHead(), 0, "")

	// 注册房间
	_ = Service2.ServiceRoom{}.Get(userModel, service, c.ClientIP(),
		Common.ClientAgentTools{}.GetDrive(c), Common.ClientAgentTools{}.GetDriveInfo(c), req.City,
	)

	// 返回token
	token := Common.Tools{}.EncodeToken(userModel.UserId, "user", service.ServiceId, 0, service.Diversion)
	Common.ApiResponse{}.Success(c, "注册成功", gin.H{"token": token})
}
