package Response

import (
	"server/App/Model/Common"
)

type DomainRes struct {
	Common.Domain
	ServiceId int    `json:"service_id"`
	Member    string `json:"member"`
	Name      string `json:"name"`
	Head      string `json:"head"`
	TimeOut   string `json:"time_out"`
}
