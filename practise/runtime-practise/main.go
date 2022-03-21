package main

import (
	"context"
	"fmt"
	"k8s.io/klog/v2"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/containerd/containerd/integration/remote/util"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

var (
	endpoint = "unix:///run/containerd/containerd.sock"
)

func main() {
	addr, dialer, err := util.GetAddressAndDialer(endpoint)
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(dialer),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(11)),
	)
	if err != nil {
		klog.Errorf("Connect remote image service %s failed: %v", addr, err)
		log.Fatalln(err)
	}

	imageClient := runtimeapi.NewImageServiceClient(conn)
	resp, err := imageClient.ImageFsInfo(ctx, &runtimeapi.ImageFsInfoRequest{})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(resp.GetImageFilesystems())
}
