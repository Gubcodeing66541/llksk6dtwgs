package Admin

import "server/App/Http/Request"

type DomainList struct {
	Type string `form:"type" json:"type" uri:"type" xml:"type"`
	Request.PageLimit
}

type DomainListLimit struct {
	Domain        string `json:"domain" uri:"domain" xml:"domain" form:"domain"`
	Username      string `json:"username" uri:"username" xml:"username" form:"username"`
	IsBindService int    `json:"is_bind_service" uri:"is_bind_service" xml:"is_bind_service" form:"is_bind_service"`
	Page          int    `json:"page" uri:"page" xml:"page" form:"page"`
	Type          string `json:"type" uri:"type" xml:"type" form:"type"`
	Offset        int    `json:"offset" uri:"offset" xml:"offset" form:"offset"`
}

type QueryById struct {
	DomainId int `form:"domain_id" json:"domain_id" uri:"domain_id" xml:"domain_id" binding:"required"`
}

type DomainDelete struct {
	DomainId int `form:"domain_id" json:"domain_id" uri:"domain_id" xml:"domain_id" binding:"required"`
}

type DomainSave struct {
	Domain string `form:"domain" json:"domain" uri:"domain" xml:"domain" binding:"required"`
	TypeEd string `form:"type" json:"type" uri:"type" xml:"type" `
}

type DomainUpdate struct {
	Id     int    `form:"id" json:"id" uri:"id" xml:"id" binding:"required"`
	Domain string `form:"domain" json:"domain" uri:"domain" xml:"domain" binding:"required"`
	Type   string `form:"type" json:"type" uri:"type" xml:"type" binding:"required"`
	Status string `form:"status" json:"status" uri:"status"  binding:"required"`
}

type DomainEnableDisable struct {
	Id     int    `form:"id" json:"id" uri:"id" xml:"id" binding:"required"`
	Status string `form:"status" json:"status" uri:"type" status:"type" binding:"required"`
}
