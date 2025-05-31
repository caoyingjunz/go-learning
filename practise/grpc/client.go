package main

import (
	"bufio"
	"context"

	"log"
	"os"
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

	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		if err := stream.Send(&pd.Request{
			Type:    "server_call",
			Payload: []byte(text),
		}); err != nil {
			log.Fatalf("Send failed: %v", err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}
