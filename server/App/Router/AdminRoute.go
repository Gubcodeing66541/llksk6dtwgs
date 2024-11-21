package Router

import (
	"github.com/gin-gonic/gin"
	"server/App/Common"
	"server/App/Http/Handel/Admin"
)

type AdminRoute struct{}

// 接口
func (AdminRoute) BindRoute(s *gin.Engine) {
	s.POST("api/admin/auth/login", Admin.Auth{}.Login)
	admin := s.Group("api/admin", ServiceMiddleWare())
	{
		admin.POST("service/create", Admin.Service{}.Create)
		admin.POST("service/renew", Admin.Service{}.Renew)
		admin.POST("service/list", Admin.Service{}.List)
		admin.POST("service/order", Admin.Service{}.Order)

		admin.POST("domain/list", Admin.Domain{}.List)
		admin.POST("domain/query_by_id", Admin.Domain{}.QueryById)
		admin.POST("domain/delete", Admin.Domain{}.Delete)
		admin.POST("domain/update", Admin.Domain{}.Update)
		admin.POST("domain/create", Admin.Domain{}.Create)
		admin.POST("domain/enable_disable", Admin.Domain{}.EnableDisable)
		admin.POST("domain/un_bind", Admin.Domain{}.UnBind)

		admin.POST("agent/list", Admin.Agent.List)
		admin.POST("agent/create", Admin.Agent.Create)
		admin.POST("agent/add_account", Admin.Agent.AddAccount)

	}

}

func AdminMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		// token
		token := c.GetHeader("token")
		var userInfo Common.UserAuthToken
		err := Common.Tools{}.DecodeToken(token, &userInfo)
		if err != nil {
			Common.ApiResponse{}.NoAuth(c, gin.H{"err": err.Error()})
			c.Abort()
			return
		}

		if userInfo.RoleType != "manager" {
			Common.ApiResponse{}.NoAuth(c, gin.H{"role_type": userInfo.RoleType})
			c.Abort()
			return
		}

		c.Set("token", token)
		c.Set("service_id", userInfo.ServiceId)
		c.Set("role_id", userInfo.RoleId)
		c.Set("role_type", userInfo.RoleType)
		c.Next()
	}
}
