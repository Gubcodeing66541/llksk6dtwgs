package Task

import (
	"fmt"
	"io/ioutil"
	"net/http"
	Common3 "server/App/Logic/Common"
	"server/App/Model/Common"
	"server/Base"
	"time"
)

type DomainStatus struct{}

func (c DomainStatus) Run() {
	fmt.Println("执行域名检测本次任务", time.Now())
	var list []Common.Domain
	Base.MysqlConn.Find(&list, "type != 'live' ")

	for _, item := range list {
		if !c.checkDomain(item.Domain) && item.Status == "enable" {
			fmt.Println("域名检测解析未生效", item.Domain)
			Base.MysqlConn.Model(&item).Updates(map[string]interface{}{"status": "no_bind_ip"})
			err := Common3.Domain{}.Bind(item.BindServiceId)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		time.Sleep(time.Second * 1)

	}

	Base.MysqlConn.Find(&list, "type != 'live' ")
	for _, item := range list {
		if c.checkDomain(item.Domain) && item.Status == "no_bind_ip" {
			fmt.Println("域名检测解析生效，恢复解析", item.Domain)
			Base.MysqlConn.Model(&item).Updates(map[string]interface{}{"status": "enable"})
			err := Common3.Domain{}.Bind(item.BindServiceId)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		time.Sleep(time.Second * 1)
	}
}

func (c DomainStatus) checkDomain(domain string) bool {
	request, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/status", domain), nil)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("无法检测域名", domain, err.Error())
		return false
	}
	page, _ := ioutil.ReadAll(resp.Body)
	return string(page) == "ok"
}
