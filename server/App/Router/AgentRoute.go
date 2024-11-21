package Router

import (
	"github.com/gin-gonic/gin"
	"server/App/Common"
	"server/App/Http/Handel/Agent"
)

type AgentRoute struct{}

func (AgentRoute) BindRoute(s *gin.Engine) {
	s.POST("api/agent/login", Agent.Agent.LoginByCode)
	s.POST("api/agent/login_by_account", Agent.Agent.LoginByAccount)
	agent := s.Group("api/agent", AgentMiddleWare())
	{
		agent.POST("/info", Agent.Agent.Info)
		agent.POST("/create_service", Agent.Agent.CreateService)
		agent.POST("/renew_service", Agent.Agent.RenewService)
		agent.POST("/recorder", Agent.Agent.LogRecorder)
		agent.POST("/list", Agent.Agent.List)
	}
}

func AgentMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		var userInfo Common.UserAuthToken
		err := Common.Tools{}.DecodeToken(token, &userInfo)
		if err != nil {
			Common.ApiResponse{}.NoAuth(c, gin.H{"err": err.Error()})
			c.Abort()
			return
		}

		if userInfo.RoleType != "agent" {
			Common.ApiResponse{}.NoAuth(c, gin.H{"err": "非代理商用户"})
			c.Abort()
			return
		}

		c.Set("service_id", userInfo.RoleId)
		c.Set("role_id", userInfo.RoleId)
		c.Set("role_type", userInfo.RoleType)
		c.Next()
	}
}
