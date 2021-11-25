package main

import (
	"context"
	"fmt"

	"github.com/caoyingjunz/client-helm/helm"
	"github.com/caoyingjunz/client-helm/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("/Users/caoyuan/.kube/config")
	if err != nil {
		panic(err)
	}
	clientSet, err := helm.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	releases, err := clientSet.AppsV1().Helms("").List(context.TODO())
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d helm relase in the cluster\n", len(releases.Items))
}
