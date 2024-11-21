package Service

import "time"

type ServiceRoomDetail struct {
	Id        int `gorm:"primary_key;AUTO_INCREMENT"`
	UserId    int `gorm:"index:user_id_idx"`
	ServiceId int `gorm:"index:service_id_idx"`
	Drive     string
	DriveInfo string `gorm:"type:varchar(1200)"`
	IP        string
	Map       string `gorm:"type:varchar(300)"`
	Mobile    string
	Wechat    string `json:"wechat"`
	Tag       string `gorm:"type:varchar(1200)"`
	IsBind    int    // 1:绑定 0:未绑定

	CreateTime time.Time
}
