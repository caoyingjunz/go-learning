package image

type Config struct {
	Default    DefaultOption    `yaml:"default"`
	Kubernetes KubernetesOption `yaml:"kubernetes"`
}

type DefaultOption struct {
	PushKubernetes bool `yaml:"push_kubernetes"`
	PushImages     bool `yaml:"push_images"`
}

type KubernetesOption struct {
	Version string `yaml:"version"`
}
