package main

import (
	"fmt"
	"golang-learning/practise/yaml-practise/module"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// refer to https://github.com/go-yaml/yaml
// https://www.jianshu.com/p/84499381a7da

func main() {
	yamlFile, err := ioutil.ReadFile("test.yaml")
	if err != nil {
		panic(err)
	}

	config := new(model.Yaml)

	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		panic(err)
	}

	fmt.Println(config.Mysql.Host)
}
