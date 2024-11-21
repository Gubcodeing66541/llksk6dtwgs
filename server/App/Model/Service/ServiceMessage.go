package Service

import "time"

type ServiceMessage struct {
	Id         int       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ServiceId  int       `json:"service_id" gorm:"index:service_id_idx"`
	MsgType    string    `json:"msg_type"` // text map image link voice video file
	MsgInfo    string    `json:"msg_info" gorm:"type:text"`
	Status     string    `json:"status"` // enable un_enable
	Type       string    `json:"type"`   // online leave quick_reply group
	CreateTime time.Time `json:"create_time"`
	Name       string    `json:"name"`
}
