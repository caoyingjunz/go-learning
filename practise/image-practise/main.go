package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"
	"sync"

	"github.com/caoyingjunz/pixiulib/exec"
	"k8s.io/klog/v2"
)

const (
	Kubeadm = "kubeadm"
	Docker  = "docker"
)

var (
	imageRepository   = flag.String("image-repository", "pixiuio", "Choose a container registry to push (default pixiuio")
	kubernetesVersion = flag.String("kubernetes-version", "", "Choose a specific Kubernetes version for the control plane")
)

type image struct {
	kubernetesVersion string
	imageRepository   string

	exec exec.Interface
}

type KubeadmVersion struct {
	ClientVersion struct {
		GitVersion string `json:"git_version"`
	} `json:"clientVersion"`
}

type KubeadmImage struct {
	Images []string `json:"images"`
}

func (img *image) Validate() error {
	if len(img.kubernetesVersion) == 0 {
		return fmt.Errorf("failed to find kubernetes version")
	}

	kubeadmVersion, err := img.getKubeadmVersion()
	if err != nil {
		return fmt.Errorf("failed to get kubeadm version: %v", err)
	}
	if kubeadmVersion != img.kubernetesVersion {
		return fmt.Errorf("kubeadm version %s not match kubernetes version %s", kubeadmVersion, img.kubernetesVersion)
	}

	return nil
}

func (img *image) getKubeadmVersion() (string, error) {
	if _, err := img.exec.LookPath(Kubeadm); err != nil {
		return "", fmt.Errorf("failed to find %s %v", Kubeadm, err)
	}

	cmd := []string{Kubeadm, "version", "-o", "json"}
	out, err := img.exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to exec kubeadm version %v %v", string(out), err)
	}

	var kubeadmVersion KubeadmVersion
	if err := json.Unmarshal(out, &kubeadmVersion); err != nil {
		return "", fmt.Errorf("failed to unmarshal kubeadm version %v", err)
	}

	return kubeadmVersion.ClientVersion.GitVersion, nil
}

func (img *image) getImages() ([]string, error) {
	cmd := []string{Kubeadm, "config", "images", "list", "--kubernetes-version", img.kubernetesVersion, "-o", "json"}
	out, err := img.exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to exec kubeadm config images list %v %v", string(out), err)
	}

	var kubeadmImage KubeadmImage
	if err := json.Unmarshal(out, &kubeadmImage); err != nil {
		return nil, fmt.Errorf("failed to unmarshal kubeadm images %v", err)
	}

	return kubeadmImage.Images, nil
}

func (img *image) doPush(imageToPush string) error {
	// real image to push
	parts := strings.Split(imageToPush, "/")
	if len(parts) < 2 {
		return fmt.Errorf("invaild image format: %s", imageToPush)
	}
	targetImage := img.imageRepository + "/" + parts[len(parts)-1]

	klog.Infof("starting pull image %s", imageToPush)
	// start pull
	cmd := []string{Docker, "pull", imageToPush}
	if _, err := img.exec.Command(cmd[0], cmd[1:]...).CombinedOutput(); err != nil {
		klog.Errorf("failed to pull %s: %v", imageToPush, err)
		return err
	}

	klog.Infof("tag %s to %s", imageToPush, targetImage)
	tagCmd := []string{Docker, "tag", imageToPush, targetImage}
	if _, err := img.exec.Command(tagCmd[0], tagCmd[1:]...).CombinedOutput(); err != nil {
		klog.Errorf("failed to tag %s to %s: %v", imageToPush, targetImage, err)
		return err
	}

	klog.Infof("starting push image %s", imageToPush)
	pushCmd := []string{Docker, "push", targetImage}
	if _, err := img.exec.Command(pushCmd[0], pushCmd[1:]...).CombinedOutput(); err != nil {
		klog.Errorf("failed to push %s: %v", imageToPush, err)
		return err
	}

	klog.Infof("complete push image %s", imageToPush)
	return nil
}

func (img *image) Push() error {
	imgs, err := img.getImages()
	if err != nil {
		return err
	}
	klog.V(2).Infof("kubeadm get images: %v", imgs)

	diff := len(imgs)
	errCh := make(chan error, diff)

	var wg sync.WaitGroup
	wg.Add(diff)
	for _, i := range imgs {
		go func(imageToPush string) {
			defer wg.Done()
			if err := img.doPush(imageToPush); err != nil {
				errCh <- err
			}
		}(i)
	}

	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	default:
	}

	return nil
}

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	img := image{
		imageRepository:   *imageRepository,
		kubernetesVersion: *kubernetesVersion,
		exec:              exec.New(),
	}

	if err := img.Validate(); err != nil {
		klog.Fatal(err)
	}

	if err := img.Push(); err != nil {
		klog.Fatal(err)
	}
}
