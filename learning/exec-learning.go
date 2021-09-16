package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func main() {

	fullArgs := make([]string, 0)
	fullArgs = append(fullArgs, []string{"-al"}...)

	output, _ := exec.Command("ls", fullArgs...).CombinedOutput()
	fmt.Println(string(output))

	cmd := exec.Command("ls", []string{"-al"}...)
	buffer := bytes.NewBuffer(nil)
	cmd.Stdout = buffer
	stderrBuffer := bytes.NewBuffer(nil)
	cmd.Stderr = stderrBuffer

	err := cmd.Run()
	if err != nil {
		fmt.Println(stderrBuffer)
	}
	fmt.Println(buffer)
}
