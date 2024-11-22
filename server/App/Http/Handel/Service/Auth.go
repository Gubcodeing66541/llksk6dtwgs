package Service

import (
	"github.com/gin-gonic/gin"
	"server/App/Common"
	Service3 "server/App/Http/Request/Service"
	Service2 "server/App/Model/Service"
	"server/Base"
	"time"
)

type Auth struct{}

// @summary 登录
// @tags 客服信息
// @Param member query string true "账号"
// @Router /api/service/auth/login [post]
func (Auth) Login(c *gin.Context) {
	var LoginReq Service3.Login
	err := c.ShouldBind(&LoginReq)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请输入正确的账号或密码", gin.H{})
		return

	}

	var ServiceModel Service2.Service
	Base.MysqlConn.Find(&ServiceModel, "member = ?", LoginReq.Member)
	if ServiceModel.ServiceId == 0 {
		Common.ApiResponse{}.Error(c, "授权账号有误", gin.H{})
		return
	}

	// 到期时间
	if ServiceModel.FirstLoginStatus == 0 {
		ServiceModel.TimeOut = time.Now().Add(time.Duration(ServiceModel.ConsumeDay) * 24 * time.Hour) // 首次登录直接天数 *24 小时重置
		ServiceModel.ConsumeDay = 0                                                                    // 重置消费天数   新开号 消费天数为 开号的天数
		ServiceModel.FirstLoginStatus = 1                                                              // 修改为非首次登录
		Base.MysqlConn.Where("service_id = ?", ServiceModel.ServiceId).Save(&ServiceModel)
	}

	// 账号到期
	if ServiceModel.TimeOut.Before(time.Now()) {
		Common.ApiResponse{}.Error(c, "账号已到期", gin.H{})
		return
	}

	token := Common.Tools{}.EncodeToken(ServiceModel.ServiceId, "service", ServiceModel.ServiceId, 0, "")
	Common.ApiResponse{}.Success(c, "ok", gin.H{"token": token})
}
