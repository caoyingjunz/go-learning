package main

import (
	"context"
	"errors"
	"fmt"
	"go-learning/practise/gorm-practise/dbstone"
)

func main() {
	udb := dbstone.NewUserDB()

	// 以使用乐观锁的方式更新 user
	updates := make(map[string]interface{})
	updates["age"] = 19
	if err := udb.OptimisticUpdate(context.TODO(), "caoyingjun", 6, updates); err != nil {
		if errors.Is(err, dbstone.RecordNotUpdated) {
			fmt.Println("hahaha not updated")
		}
	}
}
