package Common

import (
	"fmt"
	"server/App/Common"
	"server/Base"
)

type CreateMemberReq struct {
	Day uint64 `json:"day" form:"day" binding:"required"`
}

type RenewMemberReq struct {
	Member string `json:"member" form:"member" binding:"required"`
	Day    uint64 `json:"day" form:"day" binding:"required"`
}

type PostMemberReq struct {
	Member string `json:"member" form:"member" `
	Day    uint64 `json:"day" form:"day" binding:"required"`
}

type CallbackReq struct {
	Id     int64  `json:"user_order_id" form:"user_order_id" uri:"user_order_id" binding:"required"`
	Type   string `json:"type" form:"type" uri:"type" binding:"required"`
	Member string `json:"member" form:"member" uri:"member" `
	Day    uint64 `json:"day" form:"day" uri:"day" binding:"required"`
	Sign   string `json:"sign" form:"sign" uri:"sign" binding:"required"`
	Time   int64  `json:"time" form:"time" uri:"time" binding:"required"`
}

func (req *CallbackReq) AuthSign() bool {
	sign := Common.Tools{}.Md5(fmt.Sprintf("type=%s&member=%s&day=%d&key=%s&time=%d&user_order_id=%d",
		req.Type, req.Member, req.Day, Base.AppConfig.AesKey, req.Time, req.Id))
	return sign == req.Sign
}

type ShowReq struct {
	Id   int64  `json:"id" uri:"id" form:"id" binding:"required"`
	Time int64  `json:"time" uri:"time" form:"time" binding:"required"`
	Sign string `json:"sign" uri:"sign" form:"sign" binding:"required"`
}

func (req *ShowReq) CheckAuthSign() bool {
	return req.Sign == Common.Tools{}.Md5(fmt.Sprintf("%d%d%s", req.Id, req.Time, Base.AppConfig.AesKey))
}
