package main

import "fmt"

type TestInterface interface {
	RunTest1()
	RunTest2()
	RunTest3()
}

type test struct {
	*Test1
	*Test2
	*Test3
}

func NewTest() TestInterface {
	return test{
		&Test1{Name: "test1"},
		&Test2{Name: "test2"},
		&Test3{Name: "test3"},
	}
}

type Test1 struct {
	Name string
}

type Test2 struct {
	Name string
}

type Test3 struct {
	Name string
}

func (t *Test1) RunTest1() {
	fmt.Println(t.Name)
}

func (t *Test2) RunTest2() {
	fmt.Println(t.Name)
}

func (t *Test3) RunTest3() {
	fmt.Println(t.Name)
}

func main() {
	t := NewTest()

	t.RunTest1()
	t.RunTest2()
	t.RunTest3()
}
