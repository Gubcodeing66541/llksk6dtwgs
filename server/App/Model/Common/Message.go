package Common

import "time"

type Message struct {
	Id         int       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	From       int       `json:"from"`
	To         int       `json:"to"`
	Type       string    `json:"type"`
	Content    string    `json:"content" gorm:"type:text"`
	SendRole   string    `json:"send_role"` // user activate
	CreateTime time.Time `json:"create_time"`
	Diversion  string    `json:"diversion"` // 1 分流码
	IsRead     int       `json:"is_read"`   //1 已读  0未读
	UserId     int       `json:"user_id" gorm:"index:user_id_idx"`
	ServiceId  int       `json:"service_id" gorm:"index:service_id_idx"`
	IsDel      int       `json:"is_del"` //1 del
	Time       int64     `json:"time"`
}
