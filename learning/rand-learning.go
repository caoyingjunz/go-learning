package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// seed 用来确保随机，否则就是固定的值
	rand.Seed(time.Now().UnixNano())
	fmt.Println(rand.Intn(10))
	fmt.Println(rand.Perm(5))
}
