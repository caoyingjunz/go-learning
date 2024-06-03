package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("ddd")
	eventFS := make([]string, 0)

	fmt.Println(eventFS)
	fs := strings.Join(eventFS, ",")
	fmt.Println("ddd", fs, len(fs))
}
