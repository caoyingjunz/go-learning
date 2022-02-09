package main

// refer to https://github.com/spf13/cobra
import (
	"os"

	"k8s.io/component-base/cli"

	"go-learning/practise/cobra-practise/pixiuctl/app"
)

func main() {
	command := app.NewDefaultPixiuCommand()
	code := cli.Run(command)
	os.Exit(code)
}
