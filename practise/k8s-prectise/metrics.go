package main

import (
	"context"
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubernetes/pkg/controller/podautoscaler/metrics"
	"k8s.io/metrics/pkg/client/custom_metrics"
)

func main() {
	clientConfig, err := clientcmd.BuildConfigFromFlags("", "/Users/caoyuan/.kube/config")
	if err != nil {
		panic(err)
	}
	hpaClient, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		panic(err)
	}

	apiVersionsGetter := custom_metrics.NewAvailableAPIsGetter(hpaClient.Discovery())

	hpaClient.Discovery()
}
