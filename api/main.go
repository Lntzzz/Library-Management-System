package main

import (
	"Library-Management-System/api/handler"
	"log"
	"net/http"
)

func main() {
	r := handler.Init()

	// 启动 HTTP 服务器
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
