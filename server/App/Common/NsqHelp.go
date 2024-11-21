package Common

import (
	"github.com/nsqio/go-nsq"
	"server/Base"
)

type NsqHelp struct{}

func (NsqHelp) Push(topic string, msg []byte) error {
	config := nsq.NewConfig()
	p, err := nsq.NewProducer(Base.AppConfig.Mq.Nsq.Host, config)

	if err != nil {
		return err
	}

	err = p.Publish(topic, msg)
	if err != nil {
		return err
	}
	p.Stop()
	return nil
}
