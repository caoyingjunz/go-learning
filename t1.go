package main

import (
	"fmt"
)

const (
	Yellow = 100
	Red
	Green = iota * 2
	Blue
)

func T() func(string) {
	return func(s string) {
		fmt.Println(s)
	}
}

func main() {
	t := T()

	fmt.Println(t)

	t("sss")
}
