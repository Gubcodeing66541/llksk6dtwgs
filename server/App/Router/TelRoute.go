package Router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type TelRoute struct{}

func (TelRoute) BindRoute(s *gin.Engine) {
	s.StaticFS("service", http.Dir("./Tel/dist/service"))
	s.StaticFS("admin", http.Dir("./Tel/dist/admin"))
	s.StaticFS("user", http.Dir("./Tel/dist/user"))
	s.StaticFS("live", http.Dir("./Tel/dist/live"))
	s.StaticFS("common", http.Dir("./Tel/common"))
	s.StaticFS("pay", http.Dir("./Tel/dist/pay"))
	s.StaticFS("ag", http.Dir("./Tel/dist/ag"))
	s.StaticFS("fu", http.Dir("./Tel/dist/fu")) // 快速模式user

	s.LoadHTMLGlob("Tel/dist/**/*.html")

	s.StaticFS("/static", http.Dir("./static"))
}
