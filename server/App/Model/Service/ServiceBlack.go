package Service

import "time"

type ServiceBlack struct {
	Id               int       `gorm:"primary_key;AUTO_INCREMENT"`
	ServiceId        int       `gorm:"index:service_id_idx" json:"service_id"`
	Type             string    `json:"type"` //拉黑类型ip user
	UserId           int       `json:"user_id" gorm:"index:user_id_idx"`
	Ip               string    `json:"ip"`
	ServiceManagerId int       `json:"service_manager_id"`
	Day              int       `json:"day"`
	TimeOut          time.Time `json:"time_out"`
	CreateTime       time.Time `json:"create_time"`
}
