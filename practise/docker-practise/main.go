package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// TODO
var repo string = "172.16.50.151:4000"

type ImageAction interface {
	PullImage(image string) error
	PushImage(image string) error
}

type JsonStruct struct {
}

func (j *JsonStruct) Load(file string, config interface{}) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal("Load file failed")
		return
	}
	err = json.Unmarshal(data, config)
	if err != nil {
		log.Fatal("Unmarshal failed", config)
		return
	}
}

func NewJsonStruct() *JsonStruct {
	return &JsonStruct{}
}

type ImageConfig struct {
	Repo   string
	Images []string
}

type Images struct {
	input   chan struct{}
	output  chan struct{}
	counter chan struct{}
}

func runCommand(action string, image string) error {
	cmd := exec.Command("docker", action, image)
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func (i *Images) PullImage(image string) error {
	return runCommand("pull", image)
}

func (i *Images) PushImage(image string) error {
	// TODO
	newImage := strings.Join([]string{repo, image}, "/")
	cmd := exec.Command("docker", "tag", image, newImage)
	cmd.Stdout = os.Stdout
	cmd.Run()
	return runCommand("push", newImage)
}

func main() {

	JsonParse := NewJsonStruct()
	conf := ImageConfig{}
	JsonParse.Load("images.json", &conf)

	var a ImageAction = new(Images)
	// TODO
	wg := sync.WaitGroup{}

	images := conf.Images
	wg.Add(len(images))
	for _, image := range images {
		go func(image string) {
			// Decrement the counter when the goroutine completes.
			defer wg.Done()
			a.PullImage(image)
			a.PushImage(image)
		}(image)
	}
	// Wait until all goroutines completed.
	wg.Wait()
}
