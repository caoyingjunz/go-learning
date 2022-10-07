package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

var messageCh = make(chan string, 20)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.POST("/send/message", sendMessage)

	// 启动消息处理协程
	for i := 0; i < 5; i++ {
		go handleMessage()
	}

	_ = r.Run(":8080")
}

func sendMessage(c *gin.Context) {
	messageCh <- "message coming"

	c.JSON(200, "ok")
}

func handleMessage() {
	for {
		time.Sleep(5 * time.Second)
		msg := <-messageCh
		fmt.Println(msg)
	}
}
