package main

import (
	"fmt"
	"strings"
)

func main() {
	a := ""
	b := strings.Split(a, ",")
	fmt.Println(len(b))
	fmt.Println(b[1])
}
