package Common

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	Common2 "server/App/Common"
	"server/App/Http/Request/Common"
	Common3 "server/App/Logic/Common"
)

type Group struct{}

func (Group) Send(c *gin.Context) {
	var req Common.GroupSocketMessage
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, req.GetErr(err.(validator.ValidationErrors)), gin.H{})
		return
	}

	serviceId := Common2.Tools{}.GetServiceId(c)
	RoleId := Common2.Tools{}.GetRoleId(c)
	RoleType := Common2.Tools{}.GetRoleType(c)

	Common3.Group{}.SendMessage(serviceId, RoleId, RoleType, req)
	Common2.ApiResponse{}.Success(c, "ok", gin.H{})
}
