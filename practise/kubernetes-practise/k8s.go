package main

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/caoyuan/.kube/config")
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	services, err := clientset.CoreV1().Services("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, d := range services.Items {
		fmt.Println(d.Name)
	}

}
