package main

import "fmt"

func main() {

	fmt.Println("print")
	defer fmt.Println("11111111")
	panic("down down down")
	defer fmt.Println("22222222")

}
