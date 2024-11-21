package Agent

type AgentAccountLogType string

const (
	AgentAccountLogTypeRecharge      AgentAccountLogType = "recharge"       //充值
	AgentAccountLogTypeCreateService AgentAccountLogType = "create_service" //创建服务
	AgentAccountLogTypeRenewService  AgentAccountLogType = "renew_service"  //续费服务
)

type AgentAccountLog struct {
	Id            int64               `json:"id"`             // 日志ID
	AgentId       int64               `json:"agent_id"`       // 代理ID
	UpdateAccount int64               `json:"update_account"` // 变动金额
	Type          AgentAccountLogType `json:"type"`           // 类型
	ServiceCode   string              `json:"code"`           // 服务激活码
	ServiceDay    int64               `json:"day"`            // 服务天数
	ServicePrice  int64               `json:"price"`          // 服务价格
	AgentPrice    int64               `json:"agent_price"`    // 代理价格
	CreateTime    int64               `json:"create_time"`    // 创建时间
	UpdateTime    int64               `json:"update_time"`    // 更新时间
}
