package Common

import "time"

type Order struct {
	OrderId    int       `json:"order_id" gorm:"primary_key;AUTO_INCREMENT"`
	ServiceId  int       `json:"service_id"`
	Day        int       `json:"day"`
	Price      int       `json:"price"`
	Money      int       `json:"money"`
	Type       string    `json:"type"`
	CreateTime time.Time `json:"create_time"`
}
