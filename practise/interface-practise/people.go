package main

import "fmt"

type People interface {
	GetName() string
	GetAge() int

	GetData() people
}

type people struct {
	name string
}

func (p people) GetName() string {
	return p.name
}

func (p people) GetData() people {
	return p
}

func NewPeople(name string) people {
	return people{
		name: name,
	}
}

type Ming struct {
	people
	age int
}

func (x Ming) GetAge() int {
	return x.age
}

//func (x Ming) GetName() string {
// return "daming"
//}

func NewMing() People {
	return Ming{
		people: NewPeople("xiaoming"),
		age:    18,
	}
}

func main() {
	m := NewMing()

	fmt.Println(m.GetAge())
	fmt.Println(m.GetName())
	fmt.Println(m.GetData())
}
