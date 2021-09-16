package main

import "fmt"

func main() {

	stopCh := make(chan struct{})

	go func() {
		for _, i := range []string{"1", "2", "3", "4"} {
			fmt.Println(i)
		}
		close(stopCh)
	}()

	<-stopCh
}
