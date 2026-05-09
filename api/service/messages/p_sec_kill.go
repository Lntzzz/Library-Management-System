package messages

import (
	"Library-Management-System/api/internal/mq"
	"encoding/json"
	"errors"
	"fmt"
)

// 生产消息，表示资源创建结果
var _defaultSecKillProducer *SecKillProducer

func NewSecKillProducer(kConf *KafkaConf) *SecKillProducer {
	if _defaultSecKillProducer == nil {
		_defaultSecKillProducer = newSecKillProducer(kConf)
	}
	return _defaultSecKillProducer
}

func GetSecKillProducer() *SecKillProducer {
	return _defaultSecKillProducer
}

type SecKillProducer struct {
	syncProducer mq.Producer
	topic        string
}

func newSecKillProducer(kConf *KafkaConf) *SecKillProducer {
	producer := &SecKillProducer{}
	producer.topic = kConf.ProducerTopicSecKill
	producer.syncProducer = mq.NewKafkaSyncProducer(&mq.KafkaConfig{
		Name: producer.topic,
		Url:  kConf.Brokers,
	})
	return producer
}

func (p *SecKillProducer) Send(msgBody *SeckillMessage) error {

	if msgBody == nil {
		return errors.New(fmt.Sprintf("Can not send empty msg."))
	}
	content, err := json.Marshal(msgBody)
	if err != nil {
		return err
	}
	fmt.Println(string(content))
	if err := p.syncProducer.Send(p.topic, msgBody.Id, msgBody); err != nil {
		return err
	}
	return nil
}
