package main

import (
	"context"
	"fmt"

	"go-learning/practise/k8s-practise/app"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/kubernetes/pkg/apis/apps"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/metadata"
)

const (
	defaultNamespace = "default"
	defaultObject    = "nginx"
)

func main() {
	config, err := app.BuildClientConfig("")
	if err != nil {
		panic(err)
	}

	// 1. ClientSet
	// 代码地址: k8s.io/client-go/kubernetes
	// 封装 RESTClient
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	ds, err := clientSet.AppsV1().Deployments(defaultNamespace).Get(context.TODO(), defaultObject, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("clientSet", ds.Namespace, ds.Name)

	// 2. DynamicClient
	//代码地址: k8s.io/client-go/dynamic
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	unstructured, err := dynamicClient.Resource(schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}).Namespace(defaultNamespace).Get(context.TODO(), defaultObject, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	var deployment apps.Deployment
	if err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructured.UnstructuredContent(), &deployment); err != nil {
		panic(err)
	}
	fmt.Println("dynamic client", deployment.Namespace, deployment.Name)

	// 3. DiscoveryClient
	// returns the supported resources for all groups and versions.
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		panic(deployment)
	}
	_, APIResources, err := discoveryClient.ServerGroupsAndResources()
	if err != nil {
		panic(err)
	}
	fmt.Println("discovery Client", APIResources)

	// 4. metadataClient
	metadataClient, err := metadata.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	obj, err := metadataClient.Resource(schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}).Namespace(defaultNamespace).Get(context.TODO(), defaultObject, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("metadata Client", obj)

	// 5. scaleClient
	// TODO
}
