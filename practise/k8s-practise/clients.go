package main

import (
	"context"
	"fmt"

	"go-learning/practise/k8s-practise/app"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/metadata"
)

const (
	defaultNamespace = "default"
	defaultObject    = "nginx"
)

var (
	gvr = schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
)

func main() {
	config, err := app.BuildClientConfig("")
	if err != nil {
		panic(err)
	}

	// 1. ClientSet
	// 代码地址: k8s.io/client-go/kubernetes
	// 调用 k8s 的内置资源，不能访问自定义资源
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
	// 代码地址: k8s.io/client-go/dynamic
	// 通过 gvr 调用任意 k8s 资源，包括自定义资源
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	unstructured, err := dynamicClient.Resource(gvr).Namespace(defaultNamespace).Get(context.TODO(), defaultObject, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	var deployment appsv1.Deployment
	if err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructured.Object, &deployment); err != nil {
		panic(err)
	}
	fmt.Println("dynamic client", deployment.Namespace, deployment.Name)

	// 3. MetadataClient
	// 代码地址: k8s.io/client-go/metadata
	// 仅获取 k8s 对象的元数据
	// get access to a particular ObjectMeta schema without knowing the details of the version
	metadataClient, err := metadata.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	obj, err := metadataClient.Resource(gvr).Namespace(defaultNamespace).Get(context.TODO(), defaultObject, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("metadata Client", obj.Namespace, obj.Name)

	// 4. DiscoveryClient
	// 代码地址: k8s.io/client-go/discovery
	// To discover supported resources in the API server
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		panic(err)
	}
	// get the supported resources for all groups and versions.
	_, APIResources, err := discoveryClient.ServerGroupsAndResources()
	if err != nil {
		panic(err)
	}
	_ = APIResources // 忽略打印
	fmt.Println("discovery Client")

	// 5. scaleClient
	// 直接调整副本数
	// https://github.com/caoyingjunz/go-learning/blob/master/practise/k8s-practise/scale.go
}
