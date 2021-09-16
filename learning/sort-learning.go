package main

import (
	"fmt"
	"sort"
)

// 验证下 结构体排序效果

// slice 学习例子
type MyStringList []int

func (m MyStringList) Len() int {
	return len(m)
}

func (m MyStringList) Less(i, j int) bool {
	return m[i] < m[j]
}

func (m MyStringList) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

// 结构体排序
type Person struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Age  string `json:"age"`
}

type PersonSlice []Person

func (p PersonSlice) Len() int {
	return len(p)
}

func (p PersonSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PersonSlice) Less(i, j int) bool {
	return p[i].Id < p[j].Id
}

func main() {
	// 测试切片
	var names MyStringList
	names = MyStringList{3, 2, 1, 4}
	fmt.Println(names)
	sort.Sort(names)
	fmt.Println(names)

	// 测试结构体
	var persion PersonSlice
	persion = []Person{
		{
			Id:   "2",
			Name: "test2",
			Age:  "2",
		},
		{
			Id:   "1",
			Name: "test1",
			Age:  "1",
		},
		{
			Id:   "4",
			Name: "test4",
			Age:  "4",
		},
		{
			Id:   "3",
			Name: "test3",
			Age:  "3",
		},
	}

	fmt.Println(persion)
	sort.Sort(persion)
	fmt.Println(persion)
}
