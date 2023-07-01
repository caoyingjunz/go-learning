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

	IgnoreKey = "W0508"

	User     = "user"     // 修改成实际的 docker 用户名
	Password = "password" // 修改为实际的 docker 密码
)

type KubeadmVersion struct {
	ClientVersion struct {
		GitVersion string `json:"gitVersion"`
	} `json:"clientVersion"`
}

type KubeadmImage struct {
	Images []string `json:"images"`
}

type Image struct {
	KubernetesVersion string
	ImageRepository   string

	User     string
	Password string

	exec   exec.Interface
	docker *client.Client

	Cfg Config
}

func (img *Image) Validate() error {
	if img.Cfg.Default.PushKubernetes {
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
	}

	// 检查 docker 的客户端是否正常
	if _, err := img.docker.Ping(context.Background()); err != nil {
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

	if img.Cfg.Default.PushKubernetes {
		if len(img.KubernetesVersion) == 0 {
			if len(img.Cfg.Kubernetes.Version) != 0 {
				img.KubernetesVersion = img.Cfg.Kubernetes.Version
			} else {
				img.KubernetesVersion = os.Getenv("KubernetesVersion")
			}
		}
	}

	if len(img.User) == 0 {
		img.User = User
	}
	if len(img.Password) == 0 {
		img.Password = Password
	}

	img.exec = exec.New()

	if img.Cfg.Default.PushKubernetes {
		cmd := []string{"sudo", "apt-get", "install", "-y", fmt.Sprintf("kubeadm=%s-00", img.Cfg.Kubernetes.Version[1:])}
		out, err := img.exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to install kubeadm %v %v", string(out), err)
		}
	}
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

func (img *Image) cleanImages(in []byte) []byte {
	inStr := string(in)
	if !strings.Contains(inStr, IgnoreKey) {
		return in
	}

	klog.V(2).Infof("cleaning images: %+v", inStr)
	parts := strings.Split(inStr, "\n")
	index := 0
	for _, p := range parts {
		if strings.HasPrefix(p, IgnoreKey) {
			index += 1
		}
	}
	newInStr := strings.Join(parts[index:], "\n")
	klog.V(2).Infof("cleaned images: %+v", newInStr)

	return []byte(newInStr)
}

func (img *Image) getImages() ([]string, error) {
	cmd := []string{Kubeadm, "config", "images", "list", "--kubernetes-version", img.KubernetesVersion, "-o", "json"}
	out, err := img.exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to exec kubeadm config images list %v %v", string(out), err)
	}
	out = img.cleanImages(out)
	klog.V(2).Infof("images is %+v", string(out))

	var kubeadmImage KubeadmImage
	if err := json.Unmarshal(out, &kubeadmImage); err != nil {
		return nil, fmt.Errorf("failed to unmarshal kubeadm images %v", err)
	}

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

	klog.Infof("starting push image %s", targetImage)

	cmd := []string{"docker", "push", targetImage}
	out, err := img.exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to push image %s %v %v", targetImage, string(out), err)
	}

	klog.Infof("complete push image %s", imageToPush)
	return nil
}
func (img *Image) getImagesFromFile() ([]string, error) {
	var imgs []string
	for _, i := range img.Cfg.Images {
		imageStr := strings.TrimSpace(i)
		if len(imageStr) == 0 {
			continue
		}
		if strings.Contains(imageStr, " ") {
			return nil, fmt.Errorf("error image format: %s", imageStr)
		}

		imgs = append(imgs, imageStr)
	}

	return imgs, nil
}

func (img *Image) PushImages() error {
	var images []string

	if img.Cfg.Default.PushKubernetes {
		kubeImages, err := img.getImages()
		if err != nil {
			return fmt.Errorf("获取 k8s 镜像失败: %v", err)
		}
		images = append(images, kubeImages...)
	}

	if img.Cfg.Default.PushImages {
		fileImages, err := img.getImagesFromFile()
		if err != nil {
			return fmt.Errorf("")
		}
		images = append(images, fileImages...)
	}

	klog.V(2).Infof("get images: %v", images)
	diff := len(images)
	errCh := make(chan error, diff)

	// 登陆
	cmd := []string{"docker", "login", "-u", img.User, "-p", img.Password}
	out, err := img.exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to login in image %v %v", string(out), err)
	}

	var wg sync.WaitGroup
	wg.Add(diff)
	for _, i := range images {
		go func(imageToPush string) {
			defer wg.Done()
			if err := img.doPushImage(imageToPush); err != nil {
				errCh <- err
			}
		}(i)
	}
	wg.Wait()

	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	default:
	}

	return nil
}
