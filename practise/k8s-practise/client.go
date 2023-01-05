package main

import (
	"context"
	"fmt"

	"go-learning/practise/k8s-practise/app"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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

	services, err := clientSet.CoreV1().Services("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, d := range services.Items {
		fmt.Println(d.Name)
	}

}
