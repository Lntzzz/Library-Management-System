package main

import (
	"Library-Management-System/api/handler"
	"Library-Management-System/api/service/messages"
	"log"
	"net/http"
)

func main() {
	//service.Seckill.KafkaConsumerSetup()
	messages.Init()

	r := handler.Init()

	// 启动 HTTP 服务器
	err := http.ListenAndServe(":8083", r)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
