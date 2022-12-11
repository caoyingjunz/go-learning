package main

import (
	"encoding/json"
	"fmt"
)

const (
	newOperatorError = "Cannot create new Operator: %s"
)

type T1 struct {
	Name string `json:"name"`
}

func (t *T1) SetName(name string) { t.Name = name }

type T2 struct {
	T1 `json:",inline"`

	Age string `json:"age"`
}

func (t *T2) SetAge(age string) { t.Age = age }

func main() {
	in := &T2{}
	in.SetName("caoyingjunz")
	in.SetAge("19")

	d, _ := json.Marshal(in)
	fmt.Println(string(d))

	fmt.Println(fmt.Errorf(newOperatorError, "cannot create operator with nil external type"))
}
