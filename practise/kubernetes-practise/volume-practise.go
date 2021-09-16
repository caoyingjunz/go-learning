package main

import (
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

	pod, _ := clientset.CoreV1().Pods("kubez-sysns").Get("weblibriry-mariadb-0", metav1.GetOptions{})

	for _, podVolume := range pod.Spec.Volumes {
		fmt.Println(podVolume.VolumeSource)
	}
}
