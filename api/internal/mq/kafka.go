package mq

// implement Producer with kafka

import (
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
)

type saramaKafkaClient struct {
	name        string
	mu          sync.Mutex
	running     bool
	cfg         *KafkaConfig
	consumerJob chan *KafkaMessage

	consumer      *cluster.Consumer
	syncProducer  sarama.SyncProducer
	asyncProducer sarama.AsyncProducer

	noAutoMarkConsumerJob chan *sarama.ConsumerMessage
}

type KafkaMessage sarama.ConsumerMessage

type KafkaConfig struct {
	Name           string        //客户端名称，用于监控查看问题
	Url            string        //可多个，逗号分隔
	Topics         string        //可多个，逗号分隔
	GroupId        string        //消费组
	ConsumerSize   int           //并发几个消费者同时处理消息
	FromOffsets    string        //消费配置：偏移量，支持（Newest，Oldest）二种，默认使用Oldest
	CommitInterval time.Duration //消费配置：多久提交一次偏移量，默认1秒一次
	AckRule        string        //发送配置：ack规则（NoResponse、WaitForLocal、WaitForAll）默认为WaitForLocal
	AckTimeout     time.Duration //发送配置：等待Ack最大时间，默认10秒
}

func NewKafkaConsumer(cfg *KafkaConfig) Consumer {
	if cfg.GroupId == "" {
		panic(errors.New("kafka consumer 没有指定 groupId"))
	}
	config := cluster.NewConfig()
	config.Group.Return.Notifications = true
	if cfg.CommitInterval > 0 {
		config.Consumer.Offsets.CommitInterval = cfg.CommitInterval
	}
	if cfg.FromOffsets == "Newest" {
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	} else {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}
	if cfg.ConsumerSize < 1 {
		cfg.ConsumerSize = 1
	}
	consumer, err := cluster.NewConsumer(strings.Split(cfg.Url, ","), cfg.GroupId, strings.Split(cfg.Topics, ","), config)
	if err != nil {
		panic(err)
	}
	go func(c *cluster.Consumer) {
		defer OnError("")
		cerrors := c.Errors()
		noti := c.Notifications()
		for {
			select {
			case err := <-cerrors:
				if err == nil {
					return
				}
			case <-noti:
			}
		}
	}(consumer)
	return &saramaKafkaClient{
		cfg:                   cfg,
		name:                  cfg.Name,
		consumer:              consumer,
		consumerJob:           make(chan *KafkaMessage, cfg.ConsumerSize),
		noAutoMarkConsumerJob: make(chan *sarama.ConsumerMessage, cfg.ConsumerSize),
	}
}

func NewKafkaSyncProducer(cfg *KafkaConfig) Producer {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	switch cfg.AckRule {
	case "WaitForAll":
		config.Producer.RequiredAcks = sarama.WaitForAll
	case "NoResponse":
		config.Producer.RequiredAcks = sarama.NoResponse
	default:
		config.Producer.RequiredAcks = sarama.WaitForLocal
	}
	if cfg.AckTimeout > 0 {
		config.Producer.Timeout = cfg.AckTimeout
	}
	syncP, err := sarama.NewSyncProducer(strings.Split(cfg.Url, ","), config)
	if err != nil {
		panic(err)
	}
	return &saramaKafkaClient{
		cfg:          cfg,
		name:         cfg.Name,
		syncProducer: syncP,
	}
}

func NewKafkaAsyncProducer(cfg *KafkaConfig) Producer {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	switch cfg.AckRule {
	case "WaitForAll":
		config.Producer.RequiredAcks = sarama.WaitForAll
	case "NoResponse":
		config.Producer.RequiredAcks = sarama.NoResponse
	default:
		config.Producer.RequiredAcks = sarama.WaitForLocal
	}
	if cfg.AckTimeout > 0 {
		config.Producer.Timeout = cfg.AckTimeout
	}
	producer, err := sarama.NewAsyncProducer(strings.Split(cfg.Url, ","), config)
	if err != nil {
		panic(err)
	}
	go func(p sarama.AsyncProducer) {
		defer OnError("")
		perrors := p.Errors()
		success := p.Successes()
		for {
			select {
			case err := <-perrors:
				if err != nil {
					fmt.Println("KafkaAsyncProducer got error. %s", err)
				}
			case <-success:
			}
		}
	}(producer)
	return &saramaKafkaClient{
		cfg:           cfg,
		name:          cfg.Name,
		asyncProducer: producer,
	}
}

func (c *saramaKafkaClient) Send(topic, key string, values ...interface{}) error {
	for _, v := range values {
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.ByteEncoder(key),
			Value: sarama.ByteEncoder(ToByte(v)),
		}
		for i := 1; i <= 3; i++ {
			_, _, err := c.syncProducer.SendMessage(msg)
			if err == nil {
				break
			} else if i != 3 {
				continue
			} else {
				return err
			}
		}
	}
	return nil
}

func (c *saramaKafkaClient) AsyncSend(topic, key string, values ...interface{}) error {
	for _, v := range values {
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.ByteEncoder(key),
			Value: sarama.ByteEncoder(ToByte(v)),
		}
		c.asyncProducer.Input() <- msg
	}
	return nil
}

func (c *saramaKafkaClient) Listen(event func(msg interface{})) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.running {
		return
	}
	c.running = true
	// 根据并发量，启动相应个数的线程处理消息
	go func() {
		for i := 0; i < cap(c.consumerJob); i++ {
			go func(id int) {
				defer OnError("")
				for msg := range c.consumerJob {
					func() {
						defer c.onError()
						event(msg)
					}()
				}
			}(i)
		}
	}()
	// 监听消息
	go func() {
		defer OnError("")
		for msg := range c.consumer.Messages() {
			c.consumer.MarkOffset(msg, "") //MarkOffset 并不是实时写入kafka，有可能在程序crash时丢掉未提交的offset
			kafkaMsg := KafkaMessage(*msg)
			c.consumerJob <- &kafkaMsg
		}
	}()
}

func (c *saramaKafkaClient) ListenNoAutoMark(event func(msg interface{}) error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.running {
		return
	}
	c.running = true
	// 根据并发量，启动相应个数的线程处理消息
	go func() {
		for i := 0; i < cap(c.consumerJob); i++ {
			go func(id int) {
				defer OnError("")
				for msg := range c.noAutoMarkConsumerJob {
					func() {
						defer c.onError()
						kafkaMessage := KafkaMessage(*msg)
						err := event(&kafkaMessage)
						if err != nil {
						} else {
							c.consumer.MarkOffset(msg, "")
						}
					}()
				}
			}(i)
		}
	}()
	// 监听消息
	go func() {
		defer OnError("")
		for msg := range c.consumer.Messages() {
			c.noAutoMarkConsumerJob <- msg
		}
	}()
}

func (c *saramaKafkaClient) Close() {
	go func() {
		defer OnError("")
		if c.consumer != nil {
			c.consumer.Close()
		}
		if c.syncProducer != nil {
			c.syncProducer.Close()
		}
		if c.asyncProducer != nil {
			c.asyncProducer.AsyncClose()
		}
		if c.consumerJob != nil {
			close(c.consumerJob)
		}
	}()
}

func (c *saramaKafkaClient) onError() {
	if r := recover(); r != nil {
		debug.PrintStack()
	}
}

type TagResourceDelMsg struct {
	Pin         string `json:"pin"`
	Region      string `json:"region"`
	ResourceId  string `json:"resourceId"`
	ServiceCode string `json:"serviceCode"`
	Status      string `json:"status"`
}
