package Response

import "time"

type ServiceRoom struct {
	Id            int       `json:"id"`
	LateType      string    `json:"late_type"`
	UserNoRead    int       `json:"user_no_read"`
	ServiceNoRead int       `json:"service_no_read"`
	LateMsg       string    `json:"late_msg"`
	IsOnline      int       `json:"is_online"`
	IsTop         int       `json:"is_top"`
	IsBlack       int       `json:"is_black"`
	UpdateTime    time.Time `json:"update_time"`
	UserId        int       `json:"user_id"`
	UserName      string    `json:"user_name"`
	UserHead      string    `json:"user_head"`
	Rename        string    `json:"rename"`
	Times         int       `json:"times"`
}

type ServiceRoomListRes struct {
	Id            int    `json:"id"`
	LateType      string `json:"late_type"`
	UserNoRead    int    `json:"user_no_read"`
	ServiceNoRead int    `json:"service_no_read"`
	LateMsg       LateMsg
	IsOnline      int       `json:"is_online"`
	IsTop         int       `json:"is_top"`
	UpdateTime    time.Time `json:"update_time"`
	UserId        int       `json:"user_id"`
	UserName      string    `json:"user_name"`
	UserHead      string    `json:"user_head"`
	Rename        string    `json:"rename"`
}

type LateMsg struct {
	Text string `json:"text"`
	Img  string `json:"imt"`
}

type UserDetail struct {
	UserId    int    `json:"user_id"`
	UserName  string `json:"user_name"`
	UserHead  string `json:"user_head"`
	Rename    string `json:"rename"`
	IsTop     string `json:"is_top"`
	Drive     string `json:"drive"`
	DriveInfo string `json:"drive_info"`
	Ip        string `json:"ip"`
	Map       string `json:"map"`
	Mobile    string `json:"mobile"`
	Wechat    string `json:"wechat"`
	Age       string `json:"age"`
	Tag       string `json:"tag"`
	IsBind    int    `json:"is_bind"` // 是否绑定
}

type UserBlackList struct {
	Id         int       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ServiceId  int       `gorm:"index:service_id_idx" json:"service_id"`
	Type       string    `json:"type"` //拉黑类型ip user
	UserId     int       `json:"user_id"`
	Ip         string    `json:"ip"`
	UserName   string    `json:"user_name"`
	UserHead   string    `json:"user_head"`
	CreateTime time.Time `json:"create_time"`
}

type DeleteUserDay struct {
	Day int `json:"day" uri:"day" form:"day" `
}

type ServiceRoomList struct {
	UserId      int       `json:"user_id"`
	UserName    string    `json:"user_name"`
	UserHead    string    `json:"user_head"`
	Ip          string    `json:"ip"`
	Map         string    `json:"map"`
	Drive       string    `json:"drive"`
	Mobile      string    `json:"mobile"`
	Wechat      string    `json:"wechat"`
	Tag         string    `json:"tag"`
	Name        string    `json:"name"`
	IsOnline    int       `json:"is_online"`
	ServiceId   int       `json:"service_id"`
	ServiceHead string    `json:"service_head"`
	CreateTime  time.Time `json:"create_time"`
}
