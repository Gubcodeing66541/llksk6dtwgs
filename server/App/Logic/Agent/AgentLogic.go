package Agent

import (
	"errors"
	"github.com/jinzhu/gorm"
	"server/App/Common"
	Common3 "server/App/Logic/Common"
	Service22 "server/App/Logic/Service"
	"server/App/Model/Agent"
	Service2 "server/App/Model/Service"
	Setting2 "server/App/Model/Setting"
	"server/Base"
	"time"
)

type agentLogic struct{}

var AgentLogic = agentLogic{}

// CreateAgent 创建代理
func (a *agentLogic) CreateAgent(username, password string, agentPrice int64) string {
	now := time.Now()
	code := Common.Tools{}.CreateActiveCode(0)
	Base.MysqlConn.Create(&Agent.Agent{
		Username:   username,
		Password:   password,
		Code:       code,
		Status:     Agent.AgentStatusNormal,
		AgentPrice: agentPrice,
		CreateTime: now.Unix(),
		UpdateTime: now.Unix(),
	})
	return code
}

// LoginByCode 通过激活码登录
func (a *agentLogic) LoginByCode(code string) Agent.Agent {
	var agent Agent.Agent
	Base.MysqlConn.Find(&agent, "code = ?", code)
	return agent
}

// LoginByAccount 通过账号登录
func (a *agentLogic) LoginByAccount(username, password string) Agent.Agent {
	var agent Agent.Agent
	Base.MysqlConn.Find(&agent, "username = ? and password = ?", username, password)
	return agent
}

// CreateService 创建服务
func (a *agentLogic) CreateService(agentId int64, count int64, Type string, day int64) ([]string, error) {
	var agent Agent.Agent
	Base.MysqlConn.Find(&agent, "id = ?", agentId)

	agentPrice := agent.AgentPrice * count * day

	if agent.Account < agentPrice {
		return nil, errors.New("账户余额不足")
	}

	// 扣除代理账户余额
	Base.MysqlConn.Model(&agent).Update("account", agent.Account-agentPrice)

	var memberList []string
	for i := int64(0); i < count; i++ {
		member := Service22.Service{}.Create(agentId, Common.Tools{}.CreateActiveMember(), Type, int(day))
		if member.ServiceId == 0 {
			continue
		}

		memberList = append(memberList, member.Member)
		_ = Common3.Domain{}.Bind(member.ServiceId)
		a.addLog(agent, -agent.AgentPrice*day, Agent.AgentAccountLogTypeCreateService, member.Member, day)

		// 创建配置
		Base.MysqlConn.Model(&Setting2.Setting{}).Create(&Setting2.Setting{
			ServiceId:    member.ServiceId,
			Scan:         "enable",
			ScanDrive:    "default",
			ScanFitter:   "un_enable",
			ScanChange:   "un_enable",
			ToUrl:        "enable",
			Banner:       "",
			ScanToUrl:    "",
			ChangeQr:     "un_enable",
			MessageSound: "enable",
			OnlineSound:  "enable",
			Code:         member.Code,
		})
	}

	return memberList, nil
}

// RechargeAccount 充值账户
func (a *agentLogic) RechargeAccount(agentId int64, account int64) {
	var agent Agent.Agent
	Base.MysqlConn.Find(&agent, "id = ?", agentId)

	Base.MysqlConn.Model(&agent).Update("account", agent.Account+account)
	a.addLog(agent, account, Agent.AgentAccountLogTypeRecharge, "", 0)
}

// RenewService 续费服务
func (a *agentLogic) RenewService(agentId int64, memberList []string, day int64) error {
	var agent Agent.Agent
	Base.MysqlConn.Find(&agent, "id = ?", agentId)

	agentPrice := agent.AgentPrice * int64(len(memberList)) * day

	if agent.Account < agentPrice {
		return errors.New("账户余额不足")
	}

	var serviceModel Service2.Service
	var ErrMember []string
	for _, member := range memberList {
		_, err := Service22.Service{}.Renew(member, int(day), int(agent.AgentPrice))
		if err != nil {
			ErrMember = append(ErrMember, member)
			continue
		}

		// 绑定域名
		serviceModel, err = Service22.Service{}.MemberGet(member)
		if err == nil {
			_ = Common3.Domain{}.Bind(serviceModel.ServiceId)
		}
		a.addLog(agent, -agent.AgentPrice*day, Agent.AgentAccountLogTypeRenewService, member, day)
	}
	// 扣除代理账户余额

	errPrice := int64(len(memberList)) * day * agent.AgentPrice
	if agentPrice-errPrice >= 0 {
		Base.MysqlConn.Model(&agent).Update("account", gorm.Expr("account - ?", agentPrice-errPrice))
	} else {
		Base.MysqlConn.Model(&agent).Update("account", gorm.Expr("account - ?", agentPrice))
	}

	return nil
}

// addLog
func (a *agentLogic) addLog(agent Agent.Agent, UpdateAccount int64, logType Agent.AgentAccountLogType, code string, day int64) {
	Base.MysqlConn.Create(&Agent.AgentAccountLog{
		AgentId:       agent.Id,
		Type:          logType,
		ServiceDay:    day,
		ServicePrice:  agent.CreatePrice,
		ServiceCode:   code,
		UpdateAccount: UpdateAccount,
		CreateTime:    time.Now().Unix(),
	})
}

func (a *agentLogic) List(limit, offset int) []Agent.Agent {
	var list []Agent.Agent
	Base.MysqlConn.Model(&Agent.Agent{}).Limit(limit).Offset(offset).Find(&list)
	return list
}

// LogRecorder
func (a *agentLogic) LogRecorder(agentId int64, offset, limit int) []Agent.AgentAccountLog {
	var logs []Agent.AgentAccountLog
	Base.MysqlConn.Model(&Agent.AgentAccountLog{}).Where("agent_id = ?", agentId).Order("id desc").Limit(limit).Offset(offset).Find(&logs)
	return logs
}
