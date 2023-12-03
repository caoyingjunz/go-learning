package main

import (
	"context"
	"fmt"
	"go-learning/practise/k8s-practise/app"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
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

	name := "nginx"
	namespace := "default"
	deployment, err := clientSet.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	fieldSelector := []string{
		"involvedObject.uid=" + string(deployment.UID),
		"involvedObject.name=" + deployment.Name,
		"involvedObject.namespace=" + deployment.Namespace,
		"involvedObject.kind=Deployment",
	}

	events, err := clientSet.CoreV1().Events(namespace).List(context.TODO(), metav1.ListOptions{
		FieldSelector: strings.Join(fieldSelector, ","),
		Limit:         500,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("events.Items", events.Items)
	fmt.Println("events.Items", len(events.Items))
}
