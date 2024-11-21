package Service

import (
	"github.com/gin-gonic/gin"
	"server/App/Common"
	Service3 "server/App/Http/Request/Service"
	Common2 "server/App/Logic/Common"
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

	if ServiceModel.FirstLoginStatus == 0 {
		// 结束时间-创建时间
		var newTimeOut = ServiceModel.TimeOut.UnixNano() - ServiceModel.CreateTime.UnixNano()
		ServiceModel.TimeOut = time.Now().Add(time.Duration(newTimeOut))
		ServiceModel.FirstLoginStatus = 1
		Base.MysqlConn.Where("service_id = ?", ServiceModel.ServiceId).Save(&ServiceModel)
		Common2.Domain{}.Bind(ServiceModel.ServiceId)
	}

	token := Common.Tools{}.EncodeToken(ServiceModel.ServiceId, "service", ServiceModel.ServiceId, 0, "")
	Common.ApiResponse{}.Success(c, "ok", gin.H{"token": token})
}
