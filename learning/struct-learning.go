package main

import (
	"fmt"

	"golang-learning/learning/import-learning/app"
)

type People struct {
	Name string
	Age  int
}

type Property struct {
	value int
}

type myint int

func (m myint) IsZore() bool {
	return m == 0

}

func (m myint) Add(v int) int {
	return int(m) + v

}

func (p *Property) SetVaule(v int) {
	p.value = v

}

func (p *Property) GetValue() int {
	return p.value

}

func getPeople(name string, age int) *People {
	return &People{
		Name: name,
		Age:  age,
	}
}

// 实例化struct
// 1. 以 var 的方式声明结构体即可完成实例化  var p Point
// 2. 通过new来实例化  p := new(Point)
// 3. 通过 &进行实例化，p := &Point{}, 最广泛的使用方式

// 接收器是指针 修改属性
// 接收器非指针，设置无效，只能拷贝
func main() {

	//p := new(Property)
	//p.SetVaule(11)
	//fmt.Println(p.GetValue())
	//fmt.Println(p.value)

	t := &app.Test3()

	p := new(myint)

	fmt.Println(p.IsZore())
	fmt.Println(p.Add(100))
}
