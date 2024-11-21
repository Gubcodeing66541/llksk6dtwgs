package Server

import (
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
	"server/App/Common"
	"server/App/Model/Group"
	"server/Base"
	"time"
)

type NsqServer struct{}

func (NsqServer) Run(topic string) {
	cfg := nsq.NewConfig()
	cfg.LookupdPollInterval = 15 * time.Second
	c, err := nsq.NewConsumer(topic, "ch", cfg)
	if err != nil {
		panic(err)
	}
	c.AddHandler(&NsqServer{})

	//建立NSQLookupd连接
	if err := c.ConnectToNSQD(Base.AppConfig.Mq.Nsq.Host); err != nil {
		println("connect error")
		fmt.Println(Base.AppConfig.Mq.Nsq.Host, err.Error())
		if Base.AppConfig.Model != "dev" {
			panic(err)
		}
	}
}

func (NsqServer) HandleMessage(msg *nsq.Message) error {
	var req Group.GroupMessage
	err := json.Unmarshal(msg.Body, &req)
	if err != nil {
		fmt.Println("NSQ ERROR", err.Error())
	}

	// 搜索用户
	type Users struct {
		UserId int `json:"user_id"`
	}
	var users []Users
	Base.MysqlConn.Raw("select user_id from group_users where service_id = ?", req.ServiceId).Scan(&users)

	// SERVICE推送
	Common.ApiResponse{}.SendMsgToService(req.ServiceId, "group_message", req)

	// user推送
	for _, usersItem := range users {
		Common.ApiResponse{}.SendMsgToUser(usersItem.UserId, "group_message", req)
	}

	return nil
}
