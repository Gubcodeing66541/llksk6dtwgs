package User

import "time"

type User struct {
	UserId     int       `gorm:"primary_key;AUTO_INCREMENT" json:"user_id"`
	UserName   string    `json:"user_name"`
	UserHead   string    `json:"user_head"`
	OpenId     string    `json:"open_id" gorm:"index:open_id_idx"`
	UnionId    string    `json:"union_id"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	ServiceId  int       `json:"service_id"`
	Token      string    `json:"token"`
}
