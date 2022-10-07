package main

import (
	"fmt"
	apiv1 "k8s.io/api/core/v1"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println(err)
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err)
	}
	for {
		services, err := clientset.CoreV1().Services(apiv1.NamespaceDefault).List(metav1.ListOptions{})
		if err != nil {
			fmt.Println(err)
		} else {
			for _, d := range services.Items {
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Printf("%s\n", d.Name)
				}
			}
		}
	}
	time.Sleep(10 * time.Second)
}
