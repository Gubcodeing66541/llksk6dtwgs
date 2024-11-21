package Nsq

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"time"
)

type NsqConsumer struct{}

func (NsqConsumer) CreateProducer(host string) *nsq.Producer {
	// 初始化生产者
	producer, err := nsq.NewProducer(host, nsq.NewConfig())
	if err != nil {
		panic(err)
	}

	err = producer.Ping()
	if nil != err {
		// 关闭生产者
		producer.Stop()
		producer = nil
		fmt.Println("关闭生产者")
	}

	return producer
}

// 处理消息
func (NsqConsumer) HandleMessage(msg *nsq.Message) error {
	fmt.Println("receive", msg.NSQDAddress, "message:", string(msg.Body))
	return nil
}

// 初始化消费者
func (NsqConsumer) InitConsumer(topic string, channel string, host string) {
	//return
	cfg := nsq.NewConfig()
	cfg.LookupdPollInterval = time.Second * 60     //设置重连时间
	c, err := nsq.NewConsumer(topic, channel, cfg) // 新建一个消费者
	if err != nil {
		panic(err)
	}

	//c.SetLogger(nil, 0)    //屏蔽系统日志
	c.AddHandler(&NsqConsumer{}) // 添加消费者接口

	// 建立一个nsqd连接
	fmt.Println(host)
	if err := c.ConnectToNSQD(host); err != nil {
		//println(err.Error())
		panic(err)
	}
}
