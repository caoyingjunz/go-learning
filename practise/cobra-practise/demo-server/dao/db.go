package dao

import (
	"fmt"
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

type TestInterface interface {
	Get(name string) string
}

type test struct {
	factory ShareDBFactory
	db      *gorm.DB
}

func (c *test) Get(name string) string {
	fmt.Println("get name", name)
	return fmt.Sprintf("get %s", name)
}

func NewTest(f ShareDBFactory, db *gorm.DB) TestInterface {
	return &test{
		factory: f,
		db:      db,
	}
}
