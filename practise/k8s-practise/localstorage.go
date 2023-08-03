package main

import (
	"context"
	"fmt"
	"path/filepath"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/caoyingjunz/csi-driver-localstorage/pkg/client/clientset/versioned"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		panic(err)
	}
	lsClientSet, err := versioned.NewForConfig(kubeConfig)
	if err != nil {
		panic(err)
	}

	object, err := lsClientSet.StorageV1().LocalStorages().List(context.TODO(), meta_v1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println(object.Items)
}
