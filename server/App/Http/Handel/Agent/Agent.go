package Agent

import (
	"github.com/gin-gonic/gin"
	"server/App/Common"
	"server/App/Http/Request"
	Agent2 "server/App/Logic/Agent"
	Agent3 "server/App/Model/Agent"
	Service2 "server/App/Model/Service"
	"server/Base"
	"time"
)

type agent struct{}

var Agent = agent{}

// LoginByCode 通过激活码登录
func (a *agent) LoginByCode(c *gin.Context) {
	var req struct {
		Code string `json:"code" binding:"required"`
	}
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请输入登录卡密", gin.H{})
		return
	}

	agent := Agent2.AgentLogic.LoginByCode(req.Code)

	if agent.Id == 0 {
		Common.ApiResponse{}.Error(c, "卡密不存在", gin.H{})
		return
	}

	token := Common.Tools{}.EncodeToken(int(agent.Id), "agent", 0, 0, "")
	Common.ApiResponse{}.Success(c, "ok", gin.H{"token": token})
}

// LoginByAccount 通过账号登录
func (a *agent) LoginByAccount(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请输入账号密码", gin.H{})
		return
	}

	agents := Agent2.AgentLogic.LoginByAccount(req.Username, req.Password)

	if agents.Id == 0 {
		Common.ApiResponse{}.Error(c, "账号或密码错误", gin.H{})
		return
	}

	token := Common.Tools{}.EncodeToken(int(agents.Id), "agent", 0, 0, "")
	Common.ApiResponse{}.Success(c, "ok", gin.H{"token": token})
}

// CreateService 创建服务
func (a *agent) CreateService(c *gin.Context) {
	var req struct {
		Day    uint64 `json:"day" `
		Number uint64 `json:"number" binding:"required"`
		Type   string `json:"type" binding:"required"`
		Price  int64  `json:"price" `
	}
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数有误", gin.H{})
		return
	}

	agentId := Common.Tools{}.GetRoleId(c)

	code, err := Agent2.AgentLogic.CreateService(int64(agentId), int64(req.Number), req.Type, int64(req.Day))
	if err != nil {
		Common.ApiResponse{}.Error(c, err.Error(), gin.H{})
		return
	}

	Common.ApiResponse{}.Success(c, "ok", gin.H{"code": code})
}

// RenewService 续费服务
func (a *agent) RenewService(c *gin.Context) {
	var req struct {
		Services []string `json:"service" binding:"required"`
		Code     string   `json:"code" binding:"required"`
		Day      uint64   `json:"day" binding:"required"`
	}
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数有误", gin.H{})
		return
	}

	roleId := Common.Tools{}.GetRoleId(c)

	err = Agent2.AgentLogic.RenewService(int64(roleId), req.Services, int64(req.Day))
	if err != nil {
		Common.ApiResponse{}.Error(c, err.Error(), gin.H{})
		return
	}

	Common.ApiResponse{}.Success(c, "ok", gin.H{"services": req.Services})
}

// LogRecorder
func (a *agent) LogRecorder(c *gin.Context) {
	var req Request.PageLimit
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数有误", gin.H{})
		return
	}

	roleId := Common.Tools{}.GetRoleId(c)

	logs := Agent2.AgentLogic.LogRecorder(int64(roleId), req.GetOffset(), req.GetLimit())

	Common.ApiResponse{}.Success(c, "ok", gin.H{"data": logs})
}

func (a *agent) List(c *gin.Context) {
	var req Request.PageLimit
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请填写完整", gin.H{})
		return
	}
	var model []Service2.Service
	tel := Base.MysqlConn.Model(&Service2.Service{}).Where("agent_id = ?", Common.Tools{}.GetRoleId(c))
	if req.Type != "" {
		tel = tel.Where("type = ?", req.Type)
	}
	if req.Member != "" {
		tel = tel.Where("member like ?", "%"+req.Member+"%")
	}
	var Count int
	tel.Count(&Count)
	Common.DbHelp{}.ModelByPage(tel, req.Limit, req.Page).Order("service_id desc").Find(&model)
	Common.ApiResponse{}.Success(c, "ok", gin.H{"service": model, "count": Count, "current_page": req.Page})
}

func (*agent) Info(c *gin.Context) {
	roleId := Common.Tools{}.GetRoleId(c)
	var info Agent3.Agent
	Base.MysqlConn.Find(&info, "id = ?", roleId)

	var recorder []struct {
		Type string `json:"type"`
		Cnt  int64  `json:"cnt"`
	}

	now := time.Now()

	Base.MysqlConn.Raw(`
		select sum(update_account) as cnt,'today_create_service' as type from agent_account_logs
			where agent_id = ? and type = 'create_service' and create_time >= ? and create_time <= ?
			union all
			-- 今日续费
			select sum(update_account) as cnt,'today_renew_service' as type from agent_account_logs
			where agent_id = ? and type = 'renew_service' and create_time >= ? and create_time <= ?
			
			union all
			-- 今日充值
			select sum(update_account) as cnt,'today_recharge' as type from agent_account_logs
			where agent_id = ? and type = 'recharge' and create_time >= ? and create_time <= ?
			
			union all
			-- 累计开卡
			select sum(update_account) as cnt,'total_create_service' as type from agent_account_logs
			where agent_id = ? and type = 'create_service'
			
			union all
			-- 累计续费
			select sum(update_account) as cnt,'total_renew_service' as type from agent_account_logs
			where agent_id = ? and type = 'renew_service'
			
			union all
			-- 累计充值
			select sum(update_account) as cnt,'total_recharge' as type from agent_account_logs
			where agent_id = ? and type = 'recharge'
	`, roleId,
		time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Unix(),
		time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local).Unix(),
		roleId,
		time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Unix(),
		time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local).Unix(),
		roleId,
		time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Unix(),
		time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local).Unix(),
		roleId,
		roleId,
		roleId,
	).Scan(&recorder)

	Common.ApiResponse{}.Success(c, "ok", gin.H{"info": info, "recorder": recorder})
}
