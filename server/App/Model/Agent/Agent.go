package Agent

type AgentStatus string

const (
	AgentStatusNormal AgentStatus = "normal" //正常
	AgentStatusStop   AgentStatus = "stop"   //停用
)

type Agent struct {
	Id          int64       `json:"id"`           // 代理ID
	Username    string      `json:"username"`     // 代理账号
	Password    string      `json:"password"`     // 代理密码
	Code        string      `json:"code"`         // 代理激活码
	AgentPrice  int64       `json:"agent_price"`  // 代理价格
	CreatePrice int64       `json:"create_price"` // 开卡价格 代理自定义
	Trc20Addr   string      `json:"trc20_addr"`   // 代理TRC20地址
	Status      AgentStatus `json:"status"`       // 代理状态
	Ip          string      `json:"ip"`           // 代理IP
	Account     int64       `json:"account"`      // 代理账户
	CreateTime  int64       `json:"create_time"`  // 创建时间
	UpdateTime  int64       `json:"update_time"`  // 更新时间
}
