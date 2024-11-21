package User

type Token struct {
	Token string `json:"token" uri:"token" form:"token" binding:"required"`
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
	UserId     int    `json:"user_id"`
	ServiceId  int    `json:"service_id"`
}
