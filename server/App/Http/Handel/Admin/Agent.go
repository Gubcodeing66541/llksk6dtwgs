package Admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"server/App/Common"
	"server/App/Http/Request"
	Agent2 "server/App/Logic/Agent"
	"time"
)

type agent struct{}

var Agent = agent{}

func (a *agent) List(c *gin.Context) {
	var req Request.PageLimit
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请输入账号", gin.H{})
		return
	}

	agents := Agent2.AgentLogic.List(req.GetLimit(), req.GetOffset())

	Common.ApiResponse{}.Success(c, "ok", gin.H{"data": agents})
}

func (a *agent) Create(c *gin.Context) {
	var req struct {
		Username   string `json:"username" binding:"required"`
		Password   string `json:"password"`
		AgentPrice int64  `json:"agent_price" binding:"required"`
	}
	err := c.ShouldBind(&req)

	if req.Password == "" {
		req.Password = fmt.Sprintf("Pwd%d%d", time.Now().Unix(), rand.Intn(999))
	}

	if err != nil {
		Common.ApiResponse{}.Error(c, "请输入账号密码", gin.H{})
		return
	}

	code := Agent2.AgentLogic.CreateAgent(req.Username, req.Password, req.AgentPrice)

	Common.ApiResponse{}.Success(c, "ok", gin.H{"code": code})
}

func (a *agent) AddAccount(c *gin.Context) {
	var req struct {
		Code    string `json:"code" binding:"required"`
		Account int64  `json:"account" binding:"required"`
	}
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请输入卡密和金额", gin.H{})
		return
	}

	agentItem := Agent2.AgentLogic.LoginByCode(req.Code)

	if agentItem.Id == 0 {
		Common.ApiResponse{}.Error(c, "卡密不存在", gin.H{})
		return
	}

	Agent2.AgentLogic.RechargeAccount(agentItem.Id, req.Account)

	Common.ApiResponse{}.Success(c, "ok", gin.H{})
}
