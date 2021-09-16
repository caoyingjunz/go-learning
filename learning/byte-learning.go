package main

import (
	"bytes"
	"fmt"
	"sync"
)

type Porxier struct {
	mu sync.Mutex

	existingData *bytes.Buffer
}

func main() {
	proxier := &Porxier{
		existingData: bytes.NewBuffer([]byte("test")),
	}
	fmt.Println(proxier.existingData)

	proxier.existingData.Reset()
	fmt.Println("rested", proxier.existingData)

}
