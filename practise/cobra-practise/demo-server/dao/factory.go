package dao

import "github.com/jinzhu/gorm"

type ShareDBFactory interface {
	Test() TestInterface
}

type shareDBFactory struct {
	db *gorm.DB
}

func (f *shareDBFactory) Test() TestInterface {
	return NewTest(f, f.db)
}

func NewDBFactory(db *gorm.DB) ShareDBFactory {
	return &shareDBFactory{
		db: db,
	}
}
