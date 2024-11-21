package Group

import "time"

type Group struct {
	GroupId    int       `json:"group_id" gorm:"primary_key;AUTO_INCREMENT" `
	GroupName  string    `json:"group_name"`
	GroupHead  string    `json:"group_head"`
	ServiceId  int       `json:"service_id" gorm:"index:service_id_idx"`
	Status     string    `json:"status"`                     //none 无状态 stop 禁言
	Code       string    `json:"code" gorm:"index:code_idx"` //二维码
	Notice     string    `json:"notice"`                     //公告
	Hello      string    `json:"hello"`                      //进群打招呼
	LateMsg    string    `json:"late_msg"`
	NoReadCnt  int       `json:"no_read_cnt"`
	LateTime   time.Time `json:"late_time"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}
