package main

import (
	"flag"

	"k8s.io/klog/v2"

	"go-learning/practise/image-practise/image"
)

var (
	kubernetesVersion = flag.String("kubernetes-version", "", "Choose a specific Kubernetes version for the control plane")
	imageRepository   = flag.String("image-repository", "pixiuio", "Choose a container registry to push (default pixiuio")

	pushType = flag.String("type", "", "Choose the image push type")
	filePath = flag.String("file-path", "", "image file path")
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	img := image.Image{
		ImageRepository:   *imageRepository,
		KubernetesVersion: *kubernetesVersion,
		PushType:          *pushType,
		FilePath:          *filePath,
	}

	if err := img.Complete(); err != nil {
		klog.Fatal(err)
	}
	defer img.Close()

	if err := img.Validate(); err != nil {
		klog.Fatal(err)
	}

	if err := img.PushImages(); err != nil {
		klog.Fatal(err)
	}
}
