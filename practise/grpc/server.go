package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/gin-gonic/gin"
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

		s.lock.Lock()
		_, ok := s.clients[req.Type]
		if !ok {
			s.clients[req.Type] = stream
		}
		s.lock.Unlock()

		// TODO 目前是DEMO
		log.Printf("Received from %s %s", req.Type, string(req.Payload))
	}
}
func (s *server) Call(c *gin.Context) {
	_, _ = s.CallClient(c.Query("clientId"), nil)
}

func (s *server) CallClient(clientId string, data []byte) ([]byte, error) {
	stream, ok := s.clients[clientId]
	if !ok {
		return nil, fmt.Errorf("client not connected")
	}

	// 发送调用请求
	err := stream.Send(&pd.Response{
		Result: []byte(clientId + " server callback"),
	})
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

	cs := &server{clients: make(map[string]pd.Tunnel_ConnectServer)}

	s := grpc.NewServer()
	pd.RegisterTunnelServer(s, cs)

	go func() {
		log.Printf("grpc listening at %v", listener.Addr())
		if err = s.Serve(listener); err != nil {
			log.Fatalf("failed to serve %v", err)
		}
	}()

	r := gin.Default()
	r.GET("/ping", cs.Call)
	log.Printf("http listening at %v", ":8093")
	if err = r.Run(":8093"); err != nil {
		log.Fatalf("failed to start http: %v", err)
	}
}
