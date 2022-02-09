package options

import (
	"os"

	"gopkg.in/yaml.v2"

	"go-learning/practise/cobra-practise/demo-server/app/config"
)

func loadConfigFromFile(file string) (*config.Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return loadConfig(data)
}

func loadConfig(data []byte) (*config.Config, error) {
	var c config.Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	return &c, nil
}
