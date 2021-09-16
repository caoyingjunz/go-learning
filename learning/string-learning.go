package main

import (
	"fmt"
	"strings"
)

// 字符串拼接

var build strings.Builder

func main() {
	s1 := "hello"
	s2 := "world"

	// 官方推荐
	build.WriteString(s1)
	build.WriteString(" ")
	build.WriteString(s2)
	s3 := build.String()
	fmt.Println(s3)

	// pythonic
	s4 := strings.Join([]string{s1, s2}, " ")
	fmt.Println(s4)

	s5 := s1 + " " + s2
	fmt.Println(s5)

}
