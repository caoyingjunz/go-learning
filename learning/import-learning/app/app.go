package app

import "fmt"

func Test1() {
	fmt.Println("Test1")
}

func Test3() {
	fmt.Println("Test all(3)")
	Test1()
	Test2()
}
