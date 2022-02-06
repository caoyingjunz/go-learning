package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	binary, err := exec.LookPath("pixiuctl-test")
	if err != nil {
		panic(err)
	}

	args := []string{"pixiuctl-test", "test"}

	if err = syscall.Exec(binary, args, os.Environ()); err != nil {
		panic(err)
	}

	fmt.Println("end") // 不会被执行
}
