package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pd "go-learning/practise/grpc-practise/pixiu"
)

var (
	addr = "127.0.0.1:30000"
)

func main() {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect rpc server %v", err)
	}
	defer conn.Close()

	c := pd.NewPixiuClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GetPixiu(ctx, &pd.PixiuRequest{Id: 12345, Name: "caoyingjun"})
	if err != nil {
		log.Fatalf("failed to sayhello %v", err)
	}
	log.Printf("say hello %v", r.Message)
}
