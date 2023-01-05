package main

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/controller/podautoscaler/metrics"
	resourceclient "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
	"k8s.io/metrics/pkg/client/custom_metrics"
	"k8s.io/metrics/pkg/client/external_metrics"
	//autoscalingapi "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	cacheddiscovery "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/scale"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/controller-manager/pkg/clientbuilder"
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

	// 更新 scale
	//_, err = scaleClient.Scales("default").Update(context.TODO(), gr, &autoscalingapi.Scale{
	//	ObjectMeta: metav1.ObjectMeta{
	//		Name:      "test1",
	//		Namespace: "default",
	//	},
	//	Spec: autoscalingapi.ScaleSpec{
	//		Replicas: int32(4),
	//	},
	//}, metav1.UpdateOptions{})
	//if err != nil {
	//	panic(err)
	//}

	scaleKindResolver := scale.NewDiscoveryScaleKindResolver(hpaClient.Discovery())
	scaleClient, err := scale.NewForConfig(clientConfig, restMapper, dynamic.LegacyAPIPathResolverFunc, scaleKindResolver)

	// 获取 scale
	sc, err := scaleClient.Scales("default").Get(context.TODO(), gr, "test1", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println(sc.Spec.Replicas)
}
