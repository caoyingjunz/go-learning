package main

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/controller/podautoscaler/metrics"

	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/caoyuan/.kube/config")
	if err != nil {
		panic(err)
	}
	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	cs.AutoscalingV2().HorizontalPodAutoscalers("").List(context.TODO(), metav1.ListOptions{})

}
