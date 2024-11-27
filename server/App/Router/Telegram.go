package Router

import (
	"github.com/gin-gonic/gin"
	Telegram2 "server/App/Http/Handel/Telegram"
)

type Telegram struct{}

func (Telegram) BindRoute(s *gin.Engine) {
	//tg 命令
	tg := s.Group("api/tg")
	{
		tg.GET("/git/pull", Telegram2.Telegram{}.GetPull)
	}
}
