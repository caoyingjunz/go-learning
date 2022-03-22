package main

import (
	"context"
	"fmt"
	"time"

	"github.com/containerd/containerd/integration/remote/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

const (
	timeout    = 1 * time.Minute
	endpoint   = "unix:///run/containerd/containerd.sock"
	maxMsgSize = 999999
)

func main() {
	addr, dialer, err := util.GetAddressAndDialer(endpoint)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(dialer),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)),
	)
	if err != nil {
		panic(err)
	}

	imc := runtimeapi.NewImageServiceClient(conn)
	info, err := imc.ImageFsInfo(ctx, &runtimeapi.ImageFsInfoRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println(info.GetImageFilesystems())

	images, err := imc.ListImages(ctx, &runtimeapi.ListImagesRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println(images.Images)
}
