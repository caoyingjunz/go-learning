package main

import (
	cacheddiscovery "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/restmapper"
	"k8s.io/controller-manager/pkg/clientbuilder"

	"k8s.io/client-go/kubernetes"

	"go-learning/practise/k8s-practise/app"
	"k8s.io/metrics/pkg/client/custom_metrics"
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
	_ = restmapper.NewDeferredDiscoveryRESTMapper(cachedClient)

	_ = custom_metrics.NewAvailableAPIsGetter(clientSet.Discovery())

	//metricsClient := metrics.NewRESTMetricsClient(
	//	resourceclient.NewForConfigOrDie(config),
	//	custom_metrics.NewForConfig(config, mapper, apiVersionsGetter),
	//	external_metrics.NewForConfigOrDie(config),
	//)

	//fmt.Println(metricsClient)

}
