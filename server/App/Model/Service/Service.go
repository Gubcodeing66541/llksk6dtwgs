package Service

import (
	"time"
)

type Service struct {
	ServiceId        int       `json:"service_id" gorm:"primary_key;AUTO_INCREMENT"`
	AgentId          int64     `json:"agent_id"`
	Member           string    `json:"member"`
	Name             string    `json:"name"`
	Head             string    `json:"head"`
	Code             string    `json:"code" gorm:"index:code_idx"`       // 服务号
	Diversion        string    `json:"diversion"  gorm:"index:code_dis"` // 分流码
	Type             string    `json:"type"`                             // 服务类型 user group
	UserDefault      string    `json:"user_default"`                     // 用户头像是否固定头像 default 固定 空随机
	GroupCnt         int       `json:"group_cnt"`
	TimeOut          time.Time `json:"time_out"`
	CreateTime       time.Time `json:"create_time"`
	UpdateTime       time.Time `json:"update_time"`
	FirstLoginStatus int       `json:"first_login_status"` //是否第一次登录 默认0 1为登录激活
}

type ServiceMessageReply struct {
	Id        int    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ServiceId int    `json:"service_id" gorm:"index:service_id_idx"`
	Type      string `json:"type"`     // BUTTON TEXT
	MegType   string `json:"meg_type"` // image text link video
	Title     string `json:"title"`    // 标题 如果是button 则是按钮名称 如果是text 则是文本匹配内容
	Content   string `json:"content"`  // 内容
	CreateAt  int64  `json:"createAt"`
	UpdateAt  int64  `json:"updateAt"`
}
