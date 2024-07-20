package main

import (
	"flag"

	"github.com/caoyingjunz/pixiulib/config"
	"k8s.io/klog/v2"

	"go-learning/practise/image-practise/image"
)

const (
	defaultUser     = "pixiu"
	defaultPassword = "123456"
)

var (
	harbor          = flag.String("harbor", "harbor.cloud.pixiuio.com", "Choose a harbor to push (default harbor.cloud.pixiuio.com")
	imageRepository = flag.String("image-repository", "pixiuio", "Choose a container registry to push (default pixiuio")

	user     = flag.String("user", "", "docker register user")
	password = flag.String("password", "", "docker register password")

	filePath = flag.String("file-path", "", "image file path")
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	c := config.New()
	c.SetConfigFile("./config.yaml")
	c.SetConfigType("yaml")

	var cfg image.Config
	if err := c.Binding(&cfg); err != nil {
		klog.Fatal(err)
	}

	loginUser := *user
	if len(loginUser) == 0 {
		loginUser = defaultUser
	}
	loginPassword := *password
	if len(loginPassword) == 0 {
		loginPassword = defaultPassword
	}

	img := image.Image{
		Harbor:          *harbor,
		ImageRepository: *imageRepository,
		User:            loginUser,
		Password:        loginPassword,
		Cfg:             cfg,
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
