package main

import (
	"fmt"
	"reflect"
)

type S1 struct {
	Field int
}
type S2 struct {
	Field int
}

func main() {
	s1 := S1{1}
	s2 := S2{1}
	fmt.Println(reflect.DeepEqual(s1, s2))

	a := []byte{1, 3, 2}
	b := []byte{1, 3, 2}
	fmt.Println(reflect.DeepEqual(a, b))
}
