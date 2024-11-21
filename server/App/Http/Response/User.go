package Response

import "time"

type ServiceInfo struct {
	ServiceId  int       `json:"service_id"`
	Name       string    `json:"name"`
	Head       string    `json:"head"`
	Type       string    `json:"type"` //auth push
	Code       string    `json:"code"`
	Domain     string    `json:"domain"`
	Diversion  string    `json:"diversion"`
	Web        string    `json:"web"`
	TimeOut    string    `json:"time_out"`
	UserDetail string    `json:"user_detail"`
	CreateTime time.Time `json:"create_time"`
}
