package Admin

import "server/App/Http/Request"

type OrderReq struct {
	OrderId int `json:"order_id" form:"order_id" uri:"order_id" xml:"order_id"`
	Request.PageLimit
}
