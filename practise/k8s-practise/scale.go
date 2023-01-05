package main

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
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

	scaleKindResolver := scale.NewDiscoveryScaleKindResolver(clientSet.Discovery())
	scaleClient, err := scale.NewForConfig(config, restMapper, dynamic.LegacyAPIPathResolverFunc, scaleKindResolver)

	gr := schema.GroupResource{
		Group:    "apps",
		Resource: "deployments",
	}

	// 1. 获取 scale
	sc, err := scaleClient.Scales("default").Get(context.TODO(), gr, "test1", metav1.GetOptions{})
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
	fmt.Println(sc.Spec.Replicas)
}
