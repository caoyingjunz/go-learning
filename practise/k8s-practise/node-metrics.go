package main

import (
	"context"
	"fmt"

	"go-learning/practise/k8s-practise/app"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	resourceclient "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
)

func main() {
	config, err := app.BuildClientConfig("")
	if err != nil {
		panic(err)
	}

	metricClient, err := resourceclient.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	nodeMetrics, err := metricClient.NodeMetricses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, nodeMetric := range nodeMetrics.Items {
		fmt.Println(nodeMetric)
	}
}
