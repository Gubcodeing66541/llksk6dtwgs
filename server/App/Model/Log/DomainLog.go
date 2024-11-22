package Log

// CheckDomainLog 切换域名记录
type CheckDomainLog struct {
	Id         int   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ServiceId  int   `json:"service_id" gorm:"index:service_id_idx"`
	CreateTime int64 `json:"create_time"`
}
