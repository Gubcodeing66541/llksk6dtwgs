package Service

import "server/App/Http/Request"

// 续费
type Renewal struct {
	ServiceId int `json:"service_id" uri:"service_id" form:"service_id" binding:"required"`
	Day       int `json:"day" uri:"day"  form:"day" binding:"required"`
}

type GetServiceOrder struct {
	ServiceId int `json:"service_id" uri:"service_id" form:"service_id" binding:"required" `
}

type GetServiceOrderInfo struct {
	OrderId int `json:"order_id" uri:"order_id" form:"order_id" binding:"required" `
}

type BindDomain struct {
	ServiceId int `json:"service_id" uri:"service_id" form:"service_id" binding:"required" `
}

type GetServiceDomain struct {
	ServiceId int `json:"service_id" uri:"service_id" form:"service_id" binding:"required" `
}

// 创建账号天数
type CreateServiceDay struct {
	Day   int `json:"day" uri:"day"  form:"day" binding:"required"`
	Count int `json:"count" binding:"required"`
}

type GerServiceList struct {
	Request.PageLimit
	Name string `json:"name"`
}

type ChangeBindDomain struct {
	ServiceId int `json:"service_id" uri:"service_id" form:"service_id" binding:"required" `
}

type ResetQrcode struct {
}
