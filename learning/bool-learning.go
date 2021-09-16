package main

import "fmt"

func main() {
	allEndpoints := []string{"test"}

	hasEndpoints := len(allEndpoints) > 0
	fmt.Println(hasEndpoints)
}
