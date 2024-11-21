package Task

import (
	"fmt"
	"io/ioutil"
	"net/http"
	Common3 "server/App/Logic/Common"
	"server/App/Model/Common"
	"server/Base"
	"strings"
	"time"
)

type CheckDomain struct{}

func (c CheckDomain) Run() {
	time.Sleep(time.Second)
	fmt.Println("执行域名检测本次任务", time.Now())
	var list []Common.Domain
	Base.MysqlConn.Find(&list, "status = ? ", "enable")
	for _, val := range list {
		// 域名检测如果被封禁 下架域名并自动绑定已有的域名
		time.Sleep(time.Second * 3)

		status := c.checkDomain(val.Domain)
		if status == false {
			tempServiceId := val.BindServiceId
			Base.MysqlConn.Model(&val).Updates(map[string]interface{}{"bind_service_id": 0, "we_chat_ban_status": "1", "status": "un_enable"})
			err := Common3.Domain{}.Bind(tempServiceId)
			if err != nil {
				fmt.Println(err.Error())
			}

			// 推送域名封禁提示
			if val.BindServiceId != 0 {

			}
			//pararm := fmt.Sprintf("?service_id=%d&type=%s&content=%s", tempServiceId, "ban", val.Domain)
			//Common2.Tools{}.HttpGet("http://127.0.0.1/api/socket/send_to_service_socket" + pararm)

		}
		fmt.Println("check domain:", val.Domain, "status:", status)
	}
}

// https://api.urlce.com 账号yuyuyu 邮箱1106121841@qq.com 密码 6dcdW8tq
func (CheckDomain) checkDomain(domain string) bool {
	username := "yuyuyu"
	key := "hFH6vjrXRcXeoQc"
	checkUrl := fmt.Sprintf("https://api.uouin.com/app/wx?username=%s&key=%s&url=%s", username, key, domain)
	request, _ := http.NewRequest("GET", checkUrl, nil)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("无法检测域名", domain)
		return true
	}
	page, _ := ioutil.ReadAll(resp.Body)
	val := string(page)
	fmt.Println(val)
	return !(strings.Index(val, "封禁") >= 0)
}
