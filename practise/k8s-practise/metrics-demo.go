package main

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cacheddiscovery "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/controller-manager/pkg/clientbuilder"
	"k8s.io/kubernetes/pkg/controller/podautoscaler/metrics"
	resourceclient "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
	"k8s.io/metrics/pkg/client/custom_metrics"
	"k8s.io/metrics/pkg/client/external_metrics"
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

	rootClientBuilder := clientbuilder.SimpleControllerClientBuilder{
		ClientConfig: clientConfig,
	}
	discoveryClient := rootClientBuilder.DiscoveryClientOrDie("controller-discovery")
	cachedClient := cacheddiscovery.NewMemCacheClient(discoveryClient)
	restMapper := restmapper.NewDeferredDiscoveryRESTMapper(cachedClient)

	apiVersionsGetter := custom_metrics.NewAvailableAPIsGetter(hpaClient.Discovery())
	metricsClient := metrics.NewRESTMetricsClient(
		resourceclient.NewForConfigOrDie(clientConfig),
		custom_metrics.NewForConfig(clientConfig, restMapper, apiVersionsGetter),
		external_metrics.NewForConfigOrDie(clientConfig),
	)

	d, err := hpaClient.AppsV1().Deployments("default").Get(context.TODO(), "test1", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	selector, err := metav1.LabelSelectorAsSelector(d.Spec.Selector)
	// 获取 metrice cpu
	podMetricsInfo, time, err := metricsClient.GetResourceMetric(v1.ResourceMemory, "default", selector)
	if err != nil {
		panic(err)
	}

	fmt.Println("time", time)
	for podName, podMetric := range podMetricsInfo {
		fmt.Println("podName", podName)
		fmt.Println("podMetric", podMetric)
	}
}
