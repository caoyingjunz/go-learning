package image

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/caoyingjunz/pixiulib/exec"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"k8s.io/klog/v2"
)

const (
	Kubeadm = "kubeadm"
)

type KubeadmVersion struct {
	ClientVersion struct {
		GitVersion string `json:"git_version"`
	} `json:"clientVersion"`
}

type KubeadmImage struct {
	Images []string `json:"images"`
}

type Image struct {
	KubernetesVersion string
	ImageRepository   string

	exec   exec.Interface
	docker *client.Client
}

func (img *Image) Validate() error {
	if len(img.KubernetesVersion) == 0 {
		return fmt.Errorf("failed to find kubernetes version")
	}

	// 检查 kubeadm 的版本是否和 k8s 版本一致
	kubeadmVersion, err := img.getKubeadmVersion()
	if err != nil {
		return fmt.Errorf("failed to get kubeadm version: %v", err)
	}
	if kubeadmVersion != img.KubernetesVersion {
		return fmt.Errorf("kubeadm version %s not match kubernetes version %s", kubeadmVersion, img.KubernetesVersion)
	}

	// 检查 docker 的客户端是否正常
	if _, err = img.docker.Ping(context.Background()); err != nil {
		return err
	}

	return nil
}

func (img *Image) Complete() error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	img.docker = cli

	img.exec = exec.New()
	return nil
}

func (img *Image) Close() {
	if img.docker != nil {
		_ = img.docker.Close()
	}
}

func (img *Image) getKubeadmVersion() (string, error) {
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
	klog.V(2).Infof("kubeadmVersion %+v", kubeadmVersion)

	return kubeadmVersion.ClientVersion.GitVersion, nil
}

func (img *Image) getImages() ([]string, error) {
	cmd := []string{Kubeadm, "config", "images", "list", "--kubernetes-version", img.KubernetesVersion, "-o", "json"}
	out, err := img.exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to exec kubeadm config images list %v %v", string(out), err)
	}

	var kubeadmImage KubeadmImage
	if err := json.Unmarshal(out, &kubeadmImage); err != nil {
		return nil, fmt.Errorf("failed to unmarshal kubeadm images %v", err)
	}

	klog.V(2).Infof("kubeadmImage %+v", kubeadmImage)
	return kubeadmImage.Images, nil
}

func (img *Image) parseTargetImage(imageToPush string) (string, error) {
	// real image to push
	parts := strings.Split(imageToPush, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invaild image format: %s", imageToPush)
	}

	return img.ImageRepository + "/" + parts[len(parts)-1], nil
}

func (img *Image) doPushImage(imageToPush string) error {
	targetImage, err := img.parseTargetImage(imageToPush)
	if err != nil {
		return err
	}

	klog.Infof("starting pull image %s", imageToPush)
	// start pull
	reader, err := img.docker.ImagePull(context.TODO(), imageToPush, types.ImagePullOptions{})
	if err != nil {
		klog.Errorf("failed to pull %s: %v", imageToPush, err)
		return err
	}
	io.Copy(os.Stdout, reader)

	klog.Infof("tag %s to %s", imageToPush, targetImage)
	if err := img.docker.ImageTag(context.TODO(), imageToPush, targetImage); err != nil {
		klog.Errorf("failed to tag %s to %s: %v", imageToPush, targetImage, err)
		return err
	}

	klog.Infof("starting push image %s", imageToPush)
	reader, err = img.docker.ImagePush(context.TODO(), targetImage, types.ImagePushOptions{})
	if err != nil {
		klog.Errorf("failed to push %s: %v", imageToPush, err)
		return err
	}
	io.Copy(os.Stdout, reader)

	klog.Infof("complete push image %s", imageToPush)
	return nil
}

func (img *Image) PushImages() error {
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
			if err := img.doPushImage(imageToPush); err != nil {
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
