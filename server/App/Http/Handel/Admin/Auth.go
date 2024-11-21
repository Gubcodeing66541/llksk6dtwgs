package Admin

import (
	"github.com/gin-gonic/gin"
	"server/App/Common"
	"server/App/Http/Request/Admin"
	"server/Base"
)

type Auth struct{}

func (Auth) Login(c *gin.Context) {
	var req Admin.LoginReq
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "账号或密码有误", gin.H{"req": req})
		return
	}

	if req.Member != Base.AppConfig.Manager.Username || req.Password != Base.AppConfig.Manager.Password {
		Common.ApiResponse{}.Error(c, "后台账号密码有误", gin.H{})
		return
	}

	token := Common.Tools{}.EncodeToken(0, "manager", 0, 0, "")
	Common.ApiResponse{}.Success(c, "OK", gin.H{"token": token})
}
