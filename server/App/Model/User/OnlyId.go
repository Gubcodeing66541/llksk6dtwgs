package User

import "time"

type OnlyId struct {
	Id              int       `gorm:"primary_key;AUTO_INCREMENT"`
	UserId          int       `json:"user_id"`
	ServiceId       int       `json:"service_id"`
	UserToken       string    `json:"user_token" gorm:"type:text"`
	FingerprintJSID string    `json:"fingerprint_js_id"`
	CanvasID        string    `json:"canvas_id"`
	CreateTime      time.Time `json:"create_time"`
	IP              string    `json:"ip"`
}
