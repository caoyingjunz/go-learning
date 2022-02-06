package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"go-learning/practise/loadconfig-practise/config"
)

// refer to https://github.com/go-yaml/yaml

func main() {
	data, err := ioutil.ReadFile("test.yaml")
	if err != nil {
		panic(err)
	}

	config := new(config.PixiuConfiguration)
	if err = yaml.Unmarshal(data, config); err != nil {
		panic(err)
	}

	fmt.Println(config)
}
