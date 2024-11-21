package Common

type PaymentTicket struct {
	Id            int    `json:"id" gorm:"primary_key;AUTO_INCREMENT"` // 订单ID
	ServiceMember string `json:"service_member"`                       // 服务账号
	Type          string `json:"type"`                                 // create 创建 renew 续费
	Day           int    `json:"day"`                                  // 1 30 90 365
	Price         int    `json:"price"`                                // 单价 usdt
	Money         int    `json:"money"`                                // 实际支付金额
	Status        string `json:"status"`                               // wait 等待支付 success 支付成功 fail 支付失败
	CreateTime    int64  `json:"create_time"`                          // 创建时间
	UpdateTime    int64  `json:"update_time"`                          // 更新时间
	EndTime       int64  `json:"end_time"`                             // 结束时间
}
