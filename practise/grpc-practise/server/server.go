package main

// https://grpc.io/docs/languages/go/quickstart/

import (
	"context"
	pd "go-learning/practise/grpc-practise/helloworld"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	pd.UnimplementedPixiuerServer
}

func (s *server) SayHello(ctx context.Context, in *pd.HelloRequest) (*pd.HelloResponse, error) {
	log.Printf("Received %v", in.Name)
	return &pd.HelloResponse{Message: "hello " + in.GetName()}, nil
}

func main() {
	l, err := net.Listen("tcp", ":30000")
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	s := grpc.NewServer()
	pd.RegisterPixiuerServer(s, &server{})

	log.Printf("listening at %v", l.Addr())
	if err = s.Serve(l); err != nil {
		log.Fatalf("failed to serve %v", err)
	}
}
