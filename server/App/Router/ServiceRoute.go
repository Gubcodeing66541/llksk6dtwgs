package Router

import (
	"github.com/gin-gonic/gin"
	"server/App/Common"
	"server/App/Http/Handel/Service"
	"server/App/Http/Handel/Setting"
	"server/App/Http/Handel/User"
)

type ServiceRoute struct{}

func (ServiceRoute) BindRoute(s *gin.Engine) {
	s.POST("api/service/auth/login", Service.Auth{}.Login)

	s.GET("api/service/count/:code", Service.Count{}.Count)

	s.GET("api/service/get/count/:code", Service.Count{}.CountWeek)
	s.GET("api/service/get/count/month/:code", Service.Count{}.CountMonth)
	s.GET("api/service/get/count/day/:code", Service.Count{}.CountDay)

	service := s.Group("api/service", ServiceMiddleWare())
	{

		setting := service.Group("setting", ServiceMiddleWare())
		{
			setting.POST("bind_diversion", Service.Setting{}.BindDiversion)
		}

		service.POST("info", Service.Service{}.Info)
		service.POST("update", Service.Service{}.Update)

		service.POST("check_domain", Service.Service{}.CheckDomain)

		service.POST("create_live_token/:uid", Service.Service{}.CreateLive)

		// 房间相关
		service.POST("rooms/list", Service.ServiceRooms{}.List)
		service.POST("rooms/detail", Service.ServiceRooms{}.Detail)
		service.POST("rooms/update", Service.ServiceRooms{}.Update)
		service.POST("rooms/top", Service.ServiceRooms{}.Top)
		service.POST("rooms/black", Service.ServiceRooms{}.Black)
		service.POST("rooms/black_list", Service.ServiceRooms{}.BlackList)
		service.POST("rooms/rename", Service.ServiceRooms{}.Rename)
		service.POST("rooms/delete_day", Service.ServiceRooms{}.DeleteDay)
		service.POST("rooms/end", Service.ServiceRooms{}.End)

		// 消息相关
		service.POST("message/list", Service.Message{}.List)
		service.POST("message/send_all", Service.Message{}.SendAll)
		service.POST("message/send_to_user", Service.Message{}.SendToUser)
		service.POST("message/update", Service.Message{}.Update)
		service.POST("message/remove_msg", Service.Message{}.RemoveMessage)
		service.POST("message/clear_message", Service.Message{}.ClearMessage)
		service.POST("message/clear_late_message", Service.Message{}.ClearLateMessage)

		service.POST("service_message_reply/list", Service.ServiceMessageReply{}.List)

		group := service.Group("group", ServiceMiddleWare())
		{
			group.POST("detail", Service.Group{}.Detail)
			group.POST("user", Service.Group{}.User)
			group.POST("user_list", Service.Group{}.UserList)
			group.POST("update_user", Service.Group{}.UpdateUser)
			group.POST("message", User.Group{}.Message)
			group.POST("send", Service.Group{}.Send)
			group.POST("update", Service.Group{}.Update)
			group.POST("remove", Service.Group{}.Remove)
		}
	}

	// 快捷回复等相关
	serviceMessage := s.Group("api/service_message", ServiceMiddleWare())
	{
		serviceMessage.POST("create", Service.ServiceMessage{}.Create)
		serviceMessage.POST("delete", Service.ServiceMessage{}.Delete)
		serviceMessage.POST("update", Service.ServiceMessage{}.Update)
		serviceMessage.POST("list", Service.ServiceMessage{}.List)
		serviceMessage.POST("get", Service.ServiceMessage{}.GetById)
		serviceMessage.POST("swap", Service.ServiceMessage{}.Swap)
		serviceMessage.POST("set_enable", Service.ServiceMessage{}.SetEnable)

	}

	// 智能回复等相关
	botServiceMessage := s.Group("api/service_message/bot/", ServiceMiddleWare())
	{
		botServiceMessage.POST("create", Service.BotServiceMessage{}.Create)
		botServiceMessage.POST("delete", Service.BotServiceMessage{}.Delete)
		botServiceMessage.POST("update", Service.BotServiceMessage{}.Update)
		botServiceMessage.POST("list", Service.BotServiceMessage{}.List)
		botServiceMessage.POST("get", Service.BotServiceMessage{}.GetById)
		botServiceMessage.POST("swap", Service.BotServiceMessage{}.Swap)
		botServiceMessage.POST("set_enable", Service.BotServiceMessage{}.SetEnable)
	}

	// 智能回复等相关
	setting := s.Group("api/setting/", ServiceMiddleWare())
	{
		setting.POST("service_message/copy", Setting.Setting{}.CopyServiceMessage) //复制话术
		setting.POST("service_message/save", Setting.Setting{}.SaveServiceMessage) //复制话术
		setting.POST("update", Setting.Setting{}.Update)
		setting.POST("get", Setting.Setting{}.Get)
	}

}

func ServiceMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		var userInfo Common.UserAuthToken
		err := Common.Tools{}.DecodeToken(token, &userInfo)
		if err != nil {
			Common.ApiResponse{}.NoAuth(c, gin.H{"err": err.Error()})
			c.Abort()
			return
		}

		c.Set("service_id", userInfo.RoleId)
		c.Set("role_id", userInfo.RoleId)
		c.Set("role_type", userInfo.RoleType)
		c.Next()
	}
}
