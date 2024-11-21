package Common

import "time"

type UserAuthToken struct {
	RoleId    int
	RoleType  string // service user manage
	RandStr   string
	ServiceId int
	GroupId   int
	Diversion string // 分流码
	Time      time.Time
	Key       string
}
