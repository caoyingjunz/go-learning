package main

import (
	"fmt"
	"go-learning/practise/k8s-practise/app"

	"k8s.io/client-go/kubernetes"
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

	fmt.Println(clientSet)

}
