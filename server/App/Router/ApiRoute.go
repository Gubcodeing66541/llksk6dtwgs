package Router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server/App/Common"
	Common2 "server/App/Http/Handel/Common"
)

type ApiRoute struct{}

func (ApiRoute) BindRoute(s *gin.Engine) {
	s.POST("api/upload", ApiMiddleWare(), Common2.System{}.Upload)
	s.GET("websocket/conn", WebSocketMiddleWare(), Common2.WebSocketConnect{}.Conn)
	s.GET("api/live/:key", Common2.System{}.Live)
	s.POST("api/ali_upload", ApiMiddleWare(), Common2.System{}.AliUpload)

	s.GET("api/payment/callback", Common2.Payment{}.Callback)
	s.POST("api/payment/post", Common2.Payment{}.Post)
	s.POST("api/payment/return", Common2.Payment{}.Return)
	s.Any("api/payment/show/", Common2.Payment{}.Return)

	s.POST("api/app/update/ip", Common2.User.Ip)

	s.POST("test", func(context *gin.Context) {
		fmt.Println("content", context)
		context.String(200, "ok")
	})

	s.GET("api/status", func(context *gin.Context) {
		context.String(200, "ok")
	})
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

func ApiMiddleWare() gin.HandlerFunc {
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

		if userInfo.RoleType == "" {
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
