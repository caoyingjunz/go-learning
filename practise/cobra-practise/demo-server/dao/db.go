package dao

import (
	"sync"

	"github.com/jinzhu/gorm"
)

var DB *gorm.DB
var once sync.Once

func Register(db *gorm.DB) {
	once.Do(func() { // 如果是多次调用，执行 do once
		DB = db
	})
}

// TODO
