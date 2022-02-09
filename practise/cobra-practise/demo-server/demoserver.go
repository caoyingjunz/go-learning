package main

import (
	"k8s.io/klog/v2"

	"go-learning/practise/cobra-practise/demo-server/app"
)

func main() {
	cmd := app.NewDemoCommand()

	if err := cmd.Execute(); err != nil {
		klog.Fatalf("exec demo server failed %v", err)
	}
}
