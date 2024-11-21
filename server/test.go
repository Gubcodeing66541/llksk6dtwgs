package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {

	// 启动初始化
	url := "http://www.utci.top"
	checkDomains(url)

	// 启动初始化
	url = "http://www.biaoqiandayinyuezheng.com"
	checkDomains(url)

}

func checkDomains(domain string) bool {
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
