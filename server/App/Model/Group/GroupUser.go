package Group

import "time"

type GroupUser struct {
	Id         int       `json:"id" gorm:"primary_key;AUTO_INCREMENT" `
	GroupId    int       `json:"group_id" gorm:"index:group_id_idx"`
	UserId     int       `json:"user_id" gorm:"index:user_id_idx"`
	Rename     string    `json:"rename"`
	ServiceId  int       `json:"service_id" gorm:"index:service_id_idx"`
	Name       string    `json:"name"`     // 用户自己填写的名字
	Head       string    `json:"head"`     // 用户自己填写的名字
	Role       string    `json:"role"`     // user用户 manager管理员
	Status     string    `json:"status"`   // none无状态 stop禁言 ban禁止访问
	IsBlack    int       `json:"is_black"` //
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}
