package main

import (
	"fmt"
)

func TestA(a, b string, f func()) {
	fmt.Println(a)
	f()
	fmt.Println(b)
}

func main() {
	TestA("before", "end", func() {
		fmt.Println("中间执行")
	})
}
