package main

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	cacheddiscovery "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/scale"

	"go-learning/practise/k8s-practise/app"
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
	cachedClient := cacheddiscovery.NewMemCacheClient(clientSet.Discovery())

	scaleKindResolver := scale.NewDiscoveryScaleKindResolver(clientSet.Discovery())
	scaleClient, err := scale.NewForConfig(config, restmapper.NewDeferredDiscoveryRESTMapper(cachedClient), dynamic.LegacyAPIPathResolverFunc, scaleKindResolver)

	gr := schema.GroupResource{
		Group:    "apps",
		Resource: "deployments",
	}

	// 1. 获取 scale
	sc, err := scaleClient.Scales("default").Get(context.TODO(), gr, "nginx", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	newSC := sc.DeepCopy()
	newSC.Spec.Replicas = newSC.Spec.Replicas + 1
	// 更新 scale
	_, err = scaleClient.Scales("default").Update(context.TODO(), gr, newSC, metav1.UpdateOptions{})
	if err != nil {
		panic(err)
	}

	// 2. patch scale
	// TODO
}
