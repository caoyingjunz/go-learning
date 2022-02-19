package main

// refer to https://github.com/spf13/cobra
import (
	"os"

	"go-learning/practise/cobra-practise/pixiuctl/app"
)

func main() {
	cmd := app.NewDefaultPixiuCommand()

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
