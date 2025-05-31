package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"

	pd "go-learning/practise/grpc/tunnel"
)

type server struct {
	pd.UnimplementedTunnelServer

	clients map[string]pd.Tunnel_ConnectServer
	lock    sync.RWMutex
}

func (s *server) Connect(stream pd.Tunnel_ConnectServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("tream.Recv %v", err)
			return err
		}
		log.Printf("Received from %s: %s", req.Type, string(req.Payload))

		switch req.Type {
		case "client_call":
			// 处理客户端主动调用
			fmt.Println("client_call", req)
		case "server_call":
			// 这是 Server 发起的调用，Client 应返回结果
			resp := &pd.Response{Result: []byte("server pong")}
			if err = stream.Send(resp); err != nil {
				return err
			}
		}
	}
}

func (s *server) CallClient(clientID string, data []byte) ([]byte, error) {
	stream, ok := s.clients[clientID]
	if !ok {
		return nil, fmt.Errorf("client not connected")
	}

	// 发送调用请求
	err := stream.Send(&pd.Response{})
	if err != nil {
		return nil, err
	}

	return nil, err
}

func main() {
	listener, err := net.Listen("tcp", ":8092")
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	s := grpc.NewServer()
	pd.RegisterTunnelServer(s, &server{})

	log.Printf("listening at %v", listener.Addr())
	if err = s.Serve(listener); err != nil {
		log.Fatalf("failed to serve %v", err)
	}
}
