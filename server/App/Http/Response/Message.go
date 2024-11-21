package Response

import "time"

type Message struct {
	Id         int       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	From       int       `json:"from"`
	To         int       `json:"to"`
	Type       string    `json:"type"`
	Content    string    `json:"content" gorm:"type:text"`
	SendRole   string    `json:"send_role"` // user activate
	CreateTime time.Time `json:"create_time"`
	IsRead     int       `json:"is_read"`
	UserId     int       `json:"user_id" gorm:"index:idx_user_id"`
	ServiceId  int       `json:"service_id" gorm:"index:idx_service_id"`
	Time       int64     `json:"time"`
}

type SocketMessage struct {
	Id         int    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	From       int    `json:"from"`
	To         int    `json:"to"`
	Type       string `json:"type"`
	Content    string `json:"content" gorm:"type:text"`
	SendRole   string `json:"send_role"` // user activate
	CreateTime string `json:"create_time"`
	IsRead     int    `json:"is_read"`
	UserId     int    `json:"user_id" gorm:"index:idx_user_id"`
	ServiceId  int    `json:"service_id" gorm:"index:idx_service_id"`
	Time       int64  `json:"time"`
}
