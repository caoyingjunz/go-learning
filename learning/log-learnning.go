package main

import (
	"fmt"
	"log"
)

func testErr() (err error) {
	err = fmt.Errorf("test error")
	return
}

func main() {
	err := testErr()
	log.Printf("失败: %v", err)
	log.Printf("失败: %v", err.Error())
}
