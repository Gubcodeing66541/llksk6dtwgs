package Router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/App/Common"
	"server/App/Http/Handel/User"
	Common2 "server/App/Model/Common"
	"server/Base"
)

type UserRoute struct{}

func (UserRoute) BindRoute(s *gin.Engine) {

	user := s.Group("api/user")
	{
		user.GET("/long-loading", User.User{}.Loading)

		user.GET("ua/join/:code", User.UserAuth{}.Join)    // 第一步直接跳转to
		user.GET("ua/to/:code", User.UserAuth{}.To)        // 落地
		user.POST("ua/register", User.UserAuth{}.Register) // 注册

		// 全新的私聊入口
		user.GET("auth/action/:code/:uuid", User.LocalAuth{}.Action) // 第二步
		user.GET("auth/show/:code/:uuid", User.LocalAuth{}.Show)     // 第三步
		user.GET("action/:token", User.Auth{}.Action)                // 第四步 实际落地，最新的确认
		//user.GET("auth/transfer_action/:code/:uuid/:type", User.LocalAuth{}.TransferAction) //中转
		//user.GET("auth/transfer/:code/:uuid/:type", User.LocalAuth{}.Transfer)
		//中转

		// 群聊
		user.GET("auth/join_group", User.LocalAuth{}.JoinGroup)                 // 第一步 - 群聊
		user.GET("auth/action_group/:code/:uuid", User.LocalAuth{}.ActionGroup) // 第二步 - 群聊
		user.GET("auth/show_group/:code/:uuid", User.LocalAuth{}.ShowGroup)     // 第三步 - 群聊
		user.GET("action_group/detail", User.Auth{}.ActionGroup)                // 第三步 - 群聊

		user.GET("auth/join", User.LocalAuth{}.Join)                               // 私聊第一步
		user.GET("auth/diversion/action/:code/:uuid", User.AuthDiversion{}.Action) // 第二步
		user.GET("auth/location", User.LocalAuth{}.Location)                       // 中间跳转
		user.GET("auth/diversion/show/:code/:uuid", User.AuthDiversion{}.Show)     // 第三步
		user.GET("auth/diversion/:token", User.AuthDiversion{}.Html)               // 第四步 实际落地，最新的确认

		//入口
		user.POST("token", User.User{}.Token)
		user.POST("info", UserMiddleWare(), User.User{}.Info)
		user.POST("send", UserMiddleWare(), User.User{}.Send)
		user.POST("send/bot", UserMiddleWare(), User.User{}.SendBot)

		user.POST("message/list", UserMiddleWare(), User.User{}.List)
		user.POST("message/bot", UserMiddleWare(), User.User{}.BotMessage)

		user.POST("group/detail", UserMiddleWare(), User.Group{}.Detail)

		user.POST("group/user", UserMiddleWare(), User.Group{}.User)

		//user.POST("group/update", UserMiddleWare(), User.Group{}.Update)
		user.POST("group/message", UserMiddleWare(), User.Group{}.Message)
		user.POST("group/send", UserMiddleWare(), User.Group{}.Send)

		user.GET("live/create", UserMiddleWare(), User.User{}.CreateLive)

		user.GET("user/Bind", UserMiddleWare(), User.User{}.Bind)
		user.GET("user/is_bind", UserMiddleWare(), User.User{}.IsBind)

		user.POST("checkout_message", UserMiddleWare(), User.User{}.CheckoutMessage)

		user.POST("setting", UserMiddleWare(), User.User{}.Setting)

	}
}

func UserMiddleWare() gin.HandlerFunc {
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

		if userInfo.RoleType != "user" {
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

func ipCheck() gin.HandlerFunc {
	return func(context *gin.Context) {
		// token
		// 获取ip
		ip := context.ClientIP()

		// 如果IP 被禁止则直接返回
		var ipModel Common2.Ip
		Base.MysqlConn.Where("ip = ?", ip).First(&ipModel)
		if ipModel.Status == Common2.IpStatusBan {
			context.String(http.StatusForbidden, "")
			context.Abort()
			return
		}

		// 如果IP 不存在则添加
		if ipModel.Id == 0 {
			context.Redirect(http.StatusTemporaryRedirect, "/ip")
			context.Abort()
			return
		}
	}
}
