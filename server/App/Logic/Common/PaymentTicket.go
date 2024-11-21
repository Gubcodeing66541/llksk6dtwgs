package Common

import (
	"errors"
	Service22 "server/App/Logic/Service"
	"server/App/Model/Common"
	"server/Base"
	"time"
)

type PaymentTicketLogic struct{}

func (PaymentTicketLogic) CreateOrder(member string, day int, Type string, Price int, Money int) (Common.PaymentTicket, error) {
	// 搜索是否有未支付订单
	if member != "" {
		var paymentTicket Common.PaymentTicket
		Base.MysqlConn.Where("service_member = ? and status = 'wait' and day = ?", member, day).First(&paymentTicket)

		// 检查是否超过3分钟 180秒 如果超过则修改订单状态为timeout
		if paymentTicket.Id != 0 && time.Now().Unix()-paymentTicket.CreateTime > 180 {
			Base.MysqlConn.Model(&Common.PaymentTicket{}).Where("id = ?", paymentTicket.Id).
				Updates(map[string]interface{}{"status": "timeout", "update_time": time.Now().Unix()})
			paymentTicket = Common.PaymentTicket{}
		}

		if paymentTicket.Id != 0 {
			return paymentTicket, errors.New("存在未支付订单,请支付后再试，或者等待3分钟后订单过期再试")
		}
	}

	// 生成订单
	paymentTicket := Common.PaymentTicket{
		ServiceMember: member, Type: Type, Day: day, Price: Price, Money: Money, Status: "wait",
		CreateTime: time.Now().Unix(), UpdateTime: time.Now().Unix(),
	}
	Base.MysqlConn.Create(&paymentTicket)

	return paymentTicket, nil
}

func (PaymentTicketLogic) Get(id int) Common.PaymentTicket {
	var paymentTicket Common.PaymentTicket
	Base.MysqlConn.First(&paymentTicket, "id = ?", id)
	return paymentTicket
}

func (PaymentTicketLogic) Callback(id int) {
	Base.MysqlConn.Model(&Common.PaymentTicket{}).Where("id = ?", id).
		Updates(map[string]interface{}{"status": "success", "update_time": time.Now().Unix(), "end_time": time.Now().Unix()})
}

// CreateCallback 创建订单回调
func (PaymentTicketLogic) CreateCallback(Ticket Common.PaymentTicket) {
	Service22.Service{}.Create(0, Ticket.ServiceMember, "user", Ticket.Day)
}

// RenewCallback 续费订单回调
func (PaymentTicketLogic) RenewCallback(Ticket Common.PaymentTicket) {
	_, err := Service22.Service{}.Renew(Ticket.ServiceMember, Ticket.Day, Ticket.Price)
	if err != nil {
		return
	}
}
