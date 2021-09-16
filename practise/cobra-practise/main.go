package main

// refer to https://github.com/spf13/cobra

import (
	"fmt"
	"os"

	"golang-learning/practise/cobra-practise/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
