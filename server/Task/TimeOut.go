package Task

import (
	"fmt"
	"server/App/Model/Common"
	"server/Base"
	"time"
)

type TimeOut struct{}

func (TimeOut) Run() {
	fmt.Println("")
	fmt.Println("执行过期检测本次任务", time.Now())
	var list []Common.Domain
	Base.MysqlConn.Raw("select * from domains where bind_service_id in (select service_id from services where time_out <= now())").
		Scan(&list)
	for _, val := range list {
		fmt.Println("客服已过期执行下架,service_id = ", val.BindServiceId, "指定域名 = ", val.Domain)
		Base.MysqlConn.Model(&val).Updates(map[string]interface{}{"bind_service_id": 0})
	}
}
