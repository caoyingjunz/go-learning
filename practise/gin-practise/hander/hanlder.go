package hander

import (
	"context"
	"encoding/json"
	"fmt"
)

type ginMaps map[string]interface{}

type Person struct {
	Id  int64              `json:"id"`
	Age int                `json:"age"`
	Sex string             `json:"sex"`
	Max *map[string]string `json:"max"`
	Lis []string           `json:"lis"`
}

type TP struct {
	Name string `json:"name"`
}

type GinRequest struct {
	Obj  ginMaps `json:"obj"`
	Pers Person  `json:"pers"`
	TP
}

func Dohandler(ctx context.Context, gp GinRequest) error {
	// 取值
	fmt.Println("Name:", gp.Name)
	fmt.Println("Obj:", gp.Obj)

	// 序列化，可进行外部存储
	px, err := json.Marshal(gp.Pers)
	if err != nil {
		return err
	}
	fmt.Println("Marshal:", string(px))

	// 反序列化，继续使用
	var pp Person
	err = json.Unmarshal(px, &pp)
	if err != nil {
		return err
	}

	fmt.Println("Unmarshal:", pp)
	return nil
}
