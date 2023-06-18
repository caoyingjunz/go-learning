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

	obj, err := clientSet.CoreV1().Nodes().Get(context.TODO(), "kirin", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	obj.Annotations["testkey"] = "testvalue"
	_, _ = clientSet.CoreV1().Nodes().Update(context.TODO(), obj, metav1.UpdateOptions{})

	fmt.Println(obj.Annotations)
}
