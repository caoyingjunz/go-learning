package main

import "fmt"

type Test struct {
	Name string
	Age  string
}

// 接受体是指针，会保存set的属性
// 非指针，不会保存

func (t *Test) Set(name, age string) {
	t.Name = name
	t.Age = age
}

func (t *Test) GetName() string {
	return t.Name
}

func main() {

	t := &Test{
		Name: "caoyignjun",
		Age:  "18",
	}
	var x bool
	// 同时注意下defer，执行的值 只是运行到代码这行的
	t.Set("xx", "19")
	defer fmt.Println(t.GetName())
	if !x {
		t.Set("YY", "19")
	}
	fmt.Println(t.GetName())
	defer fmt.Println(t.GetName())
}
