package main

import (
	"context"
	"fmt"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/metadata"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		panic(err)
	}
	metadataClient, err := metadata.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	obj, err := metadataClient.Resource(schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}).Namespace("default").Get(context.TODO(), "nginx", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println(obj)
}
