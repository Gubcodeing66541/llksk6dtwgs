package App

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"server/App/Common"
	"server/App/Router"
)

type HttpRoute struct{}

// 限流器
var rateLimit = rate.NewLimiter(600, 800)

func (HttpRoute) BindRoute(s *gin.Engine) {
	s.Use(Cors())

	//s.Use(LimitMiddleWare())

	s.NoRoute(func(c *gin.Context) {
		// 实现内部重定向
		c.HTML(http.StatusOK, "tp.html", gin.H{})
	})

	Router.TelRoute{}.BindRoute(s)
	Router.ApiRoute{}.BindRoute(s)
	Router.AdminRoute{}.BindRoute(s)
	Router.ServiceRoute{}.BindRoute(s)
	Router.UserRoute{}.BindRoute(s)
	Router.AgentRoute{}.BindRoute(s)
}

// 所有接口进行限流
func LimitMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		if rateLimit.Allow() {
			c.Next()
		} else {
			c.String(http.StatusInternalServerError, "服务器繁忙，请稍后在试")
			c.Abort()
		}
	}
}

func WebSocketMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.DefaultQuery("token", "")
		var userInfo Common.UserAuthToken
		err := Common.Tools{}.DecodeToken(token, &userInfo)
		if err != nil {
			Common.ApiResponse{}.NoAuth(c, gin.H{"err": err.Error(), "token": token, "type": "websocket"})
			c.Abort()
			return
		}

		c.Set("service_id", userInfo.ServiceId)
		c.Set("role_id", userInfo.RoleId)
		c.Set("role_type", userInfo.RoleType)
		c.Set("group_id", userInfo.GroupId)
		c.Next()
	}
}

func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,X-Requested-With,XMLHttpRequest")
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, Token ,token")
		context.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusOK)
		}
		context.Next()
	}
}
