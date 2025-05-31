package main

import (
	"context"

	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pd "go-learning/practise/grpc/tunnel"
)

var (
	addr = "127.0.0.1:8092"
)

func main() {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect rpc server %v", err)
	}
	defer conn.Close()

	c := pd.NewTunnelClient(conn)

	stream, err := c.Connect(context.Background())
	if err != nil {
		log.Fatalf("%v", err)
	}

	clientId := "node2"

	// 启动协程，接受服务段回调 client 的请求
	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				log.Printf("Receive error: %v", err)
				return
			}
			log.Printf("Received from server: %s", msg.Result)
		}
	}()

	// 启动客户端定时测试DEMO
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		ts := time.Now().String()
		if err = stream.Send(&pd.Request{
			Type:    clientId,
			Payload: []byte(ts),
		}); err != nil {
			log.Println("调用服务端失败", err)
		}
	}
}
