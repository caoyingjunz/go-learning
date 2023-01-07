package app

import (
	"path/filepath"

	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func BuildClientConfig(configFile string) (*restclient.Config, error) {
	if len(configFile) == 0 {
		configFile = filepath.Join(homedir.HomeDir(), ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", configFile)
	if err != nil {
		return nil, err
	}
	return config, nil
}
