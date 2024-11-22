package Service

import "server/App/Http/Request"

type CreateServiceMessage struct {
	ServiceId int    `json:"service_id" uri:"service_id" form:"service_id" `
	MsgInfo   string `json:"msg_info" uri:"msg_info" form:"msg_info" `
	MsgType   string `json:"msg_type" uri:"msg_type" form:"msg_type" binding:"required" `
	Type      string `json:"type" uri:"type" form:"type" binding:"required" `
	Name      string `json:"name" uri:"name" form:"name"  `
}

type ListServiceMessage struct {
	Request.PageLimit
}

type UpdateServiceMessage struct {
	Id        int    `json:"id" uri:"id" form:"id" binding:"required" `
	MsgInfo   string `json:"msg_info" uri:"msg_info" form:"msg_info" `
	MsgType   string `json:"msg_type" uri:"msg_type" form:"msg_type" binding:"required" `
	ServiceId int    `json:"service_id" uri:"service_id" form:"service_id"  `
	Name      string `json:"name" uri:"name" form:"name"  `
	Status    string `json:"status" uri:"status" form:"status"  `
}

type DeleteServiceMessage struct {
	Id        int `json:"id" uri:"id" form:"id" binding:"required" `
	ServiceId int `json:"service_id" uri:"service_id" form:"service_id" `
}

type GetByIdServiceMessage struct {
	Id        int `json:"id" uri:"id" form:"id" binding:"required" `
	ServiceId int `json:"service_id" uri:"service_id" form:"service_id" `
}

type SwapServiceMessage struct {
	From int `json:"from" uri:"from" form:"from" binding:"required" `
	To   int `json:"to" uri:"to" form:"to" binding:"required" `
}
