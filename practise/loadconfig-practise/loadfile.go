package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"

	"go-learning/practise/loadconfig-practise/config"
)

func loadConfigFromFile(file string) (*config.PixiuConfiguration, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return loadConfig(data)
}

func loadConfig(data []byte) (*config.PixiuConfiguration, error) {
	var pc config.PixiuConfiguration
	if err := yaml.Unmarshal(data, &pc); err != nil {
		return nil, err
	}

	return &pc, nil
}

func main() {
	pixiuConfiguration, err := loadConfigFromFile("test.yaml")
	if err != nil {
		panic(err)
	}

	fmt.Println(pixiuConfiguration)
}
