package User

import (
	"time"
)

type UserLoginLog struct {
	Id         int    `gorm:"primary_key;AUTO_INCREMENT"`
	UserId     int    `json:"user_id"`
	ServiceId  int    `json:"service_id"`
	Ip         string `json:"ip"`
	Addr       string `json:"addr"`
	CreateTime time.Time
}
