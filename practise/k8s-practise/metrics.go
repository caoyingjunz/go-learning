package main

import (
	"context"
	"fmt"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cacheddiscovery "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
	"k8s.io/controller-manager/pkg/clientbuilder"
	resourceclient "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"

	"k8s.io/metrics/pkg/client/custom_metrics"
	"k8s.io/metrics/pkg/client/external_metrics"

	"go-learning/practise/k8s-practise/app"
	"go-learning/practise/k8s-practise/metrics"
)

func main() {
	config, err := app.BuildClientConfig("")
	if err != nil {
		panic(err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	rootClientBuilder := clientbuilder.SimpleControllerClientBuilder{
		ClientConfig: config,
	}
	discoveryClient := rootClientBuilder.DiscoveryClientOrDie("controller-discovery")
	cachedClient := cacheddiscovery.NewMemCacheClient(discoveryClient)
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(cachedClient)

	apiVersionsGetter := custom_metrics.NewAvailableAPIsGetter(clientSet.Discovery())

	metricsClient := metrics.NewRESTMetricsClient(
		resourceclient.NewForConfigOrDie(config),
		custom_metrics.NewForConfig(config, mapper, apiVersionsGetter),
		external_metrics.NewForConfigOrDie(config),
	)

	deploy, err := clientSet.AppsV1().Deployments("default").Get(context.TODO(), "nginx", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	selector, _ := metav1.LabelSelectorAsSelector(deploy.Spec.Selector)
	// 获取 metrice cpu
	podMetricsInfo, time, err := metricsClient.GetResourceMetric(v1.ResourceCPU, "default", selector)
	if err != nil {
		panic(err)
	}

	fmt.Println("time", time)
	for name, metric := range podMetricsInfo {
		fmt.Println("name", name, "metric", metric)
	}
}
