package main

// https://grpc.io/docs/languages/go/quickstart/

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pd "go-learning/practise/grpc-practise/pixiu"
)

type server struct {
	pd.UnimplementedPixiuServer
}

func (s *server) GetPixiu(ctx context.Context, in *pd.PixiuRequest) (*pd.PixiuReply, error) {
	log.Printf("Received %s %d", in.Name, in.Id)
	return &pd.PixiuReply{Message: fmt.Sprintf("%s %d", in.GetName(), in.GetId())}, nil
}

func main() {
	l, err := net.Listen("tcp", ":30000")
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	s := grpc.NewServer()
	pd.RegisterPixiuServer(s, &server{})

	log.Printf("listening at %v", l.Addr())
	if err = s.Serve(l); err != nil {
		log.Fatalf("failed to serve %v", err)
	}
}
