package Common

type IpStatus string

var (
	IpStatusPass IpStatus = "pass" // 正常
	IpStatusBan  IpStatus = "ban"  // 禁用
)

type Ip struct {
	Id         int      `json:"id" gorm:"primaryKey"`
	Ip         string   `json:"ip" gorm:"index:idx_ip"`
	Ext        string   `json:"ext" gorm:"type:text"`
	Status     IpStatus `json:"status"`      // 状态 0:正常 1:禁用
	CreateTime int64    `json:"create_time"` // 创建时间
}
