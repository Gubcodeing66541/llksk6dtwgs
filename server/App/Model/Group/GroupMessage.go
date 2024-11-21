package Group

import "time"

type GroupMessage struct {
	Id         int       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	GroupId    int       `json:"group_id" gorm:"index:group_id_idx"`
	ServiceId  int       `json:"service_id" gorm:"index:service_id_idx"`
	SendRoleId int       `json:"send_role_id"`
	SendRole   string    `json:"send_role"` // user service
	Type       string    `json:"type"`
	Content    string    `json:"content" gorm:"type:text"`
	ReadCnt    int       `json:"read_cnt"`
	CreateTime time.Time `json:"create_time"`
	Time       int64     `json:"time"`
}
