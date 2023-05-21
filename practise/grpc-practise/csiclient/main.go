package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/kubernetes-csi/csi-lib-utils/connection"
	"github.com/kubernetes-csi/csi-lib-utils/metrics"
	csirpc "github.com/kubernetes-csi/csi-lib-utils/rpc"
	"k8s.io/klog/v2"
)

var (
	csiAddress       = flag.String("csi-address", "/tmp/csi.sock", "Path of the CSI driver socket that the node-driver-registrar will connect to.")
	operationTimeout = flag.Duration("timeout", time.Second, "Timeout for waiting for communication with driver")
)

func main() {
	flag.Parse()

	// Unused metrics manager, necessary for connection.Connect below
	cmm := metrics.NewCSIMetricsManagerForSidecar("")

	klog.V(0).Infof("Attempting to open a gRPC connection with: %q", *csiAddress)
	csiConn, err := connection.Connect(*csiAddress, cmm)
	if err != nil {
		klog.Errorf("error connecting to CSI driver: %v", err)
		os.Exit(1)
	}

	klog.V(0).Infof("Calling CSI driver to discover driver name")
	ctx, cancel := context.WithTimeout(context.Background(), *operationTimeout)
	defer cancel()

	// identityserver rpc
	csiDriverName, err := csirpc.GetDriverName(ctx, csiConn)
	if err != nil {
		klog.Errorf("error retreiving CSI driver name: %v", err)
		os.Exit(1)
	}
	fmt.Println("csiDriverName", csiDriverName)

	ready, err := csirpc.Probe(ctx, csiConn)
	if err != nil {
		klog.Errorf("error retreiving CSI Probe: %v", err)
		os.Exit(1)
	}
	fmt.Println("ready", ready)

	// nodeserver
	nodeClient := csi.NewNodeClient(csiConn)
	nodeInfo, err := nodeClient.NodeGetInfo(ctx, &csi.NodeGetInfoRequest{})
	if err != nil {
		klog.Errorf("error NodeGetInfo: %v", err)
		os.Exit(1)
	}
	fmt.Println("NodeGetInfo", nodeInfo)

	// controllerserver rpc
	csiClient := csi.NewControllerClient(csiConn)
	resp, err := csiClient.CreateVolume(ctx, &csi.CreateVolumeRequest{
		Name:               "test-volume",
		VolumeCapabilities: []*csi.VolumeCapability{},
	})
	if err != nil {
		klog.Errorf("error CreateVolume: %v", err)
		os.Exit(1)
	}
	fmt.Println("create volume", resp)
}
