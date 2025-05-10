package mq

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"time"
)

type Starter interface {
	Start()
}
type Producer interface {

	// Send 同步发送
	Send(topic, key string, values ...interface{}) error

	// AsyncSend 异步发送
	AsyncSend(topic, key string, values ...interface{}) error

	// Close 关闭连接
	Close()
}

type Consumer interface {

	// Listen 开启消费监听
	Listen(event func(msg interface{}))

	// Close 关闭连接
	Close()
}

type NoAutoMarkConsumer interface {
	Consumer
	ListenNoAutoMark(event func(msg interface{}) error)
}

func OnError(id string) {
	if r := recover(); r != nil {
		fmt.Println(id, time.Now())
		debug.PrintStack()
	}
}

func ToByte(v interface{}) []byte {
	switch t := v.(type) {
	case string:
		return []byte(t)
	}
	j, _ := json.Marshal(v)
	return j
}
