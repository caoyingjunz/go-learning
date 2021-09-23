package main

import (
	"context"
	"fmt"
	"time"

	"go-learning/practise/gorm-practise/dbstone"
	"go-learning/practise/gorm-practise/models"
)

const (
	name = "caoyingjun"
	age  = 18
)

var udb = dbstone.NewUserDB()

func main() {
	var err error

	//  创建 user
	now := time.Now()
	u1 := models.User{
		GmtCreate:   now,
		GmtModified: now,
		Name:        name,
		Age:         age,
	}
	if err = udb.Create(context.TODO(), &u1); err != nil {
		// handler err
		panic(err)
	}

	// 获取 user
	u2, err := udb.Get(context.TODO(), name, age)
	if err != nil {
		panic(err)
	}
	fmt.Println("get user: ", u2)

	// 更新 user
	// 业务层面的数据更新
	// 数据库层面的数据在 update 中实现
	updates := make(map[string]interface{})
	updates["age"] = 19
	if err = udb.Update(context.TODO(), name, u2.ResourceVersion, updates); err != nil {
		panic(err)
	}

	users, err := udb.ListByPage(context.TODO(), name, 2, 2)
	if err != nil {
		panic(err)
	}
	fmt.Println("list page users:", users)

	// 删除 user
	//if err = udb.Delete(context.TODO(), u2.ID); err != nil {
	//	panic(err)
	//}
}
