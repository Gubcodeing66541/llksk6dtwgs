package main

import (
	"server/Base"
	"server/Task"
	"time"
)

func main() {
	// 启动初始化
	Base.Base{}.Init()

	go func() {
		for true {
			// 跑域名检测
			Task.CheckDomain{}.Run()
		}
	}()

	for {
		// 过期检测
		go Task.TimeOut{}.Run()

		// 解析失效解析
		go Task.DomainStatus{}.Run()

		// 休眠30秒
		time.Sleep(time.Second * 30)
	}
}
