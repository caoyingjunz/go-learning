package models

import "time"

//CREATE TABLE `user` (
//	`id` INT NOT NULL AUTO_INCREMENT,
//	`name` VARCHAR(16) NOT NULL,
//	`age` TINYINT NOT NULL,
//	PRIMARY KEY (`id`)
//) ENGINE=InnoDB CHARSET=utf8 AUTO_INCREMENT=20110000;

type User struct {
	ID              int64     `gorm:"column:id;primary_key;AUTO_INCREMENT;not null" json:"id"`
	GmtCreate       time.Time `json:"gmt_create"`
	GmtModified     time.Time `json:"gmt_modified"`
	ResourceVersion int64     `json:"resource_version"`
	Name            string    `json:"name"`
	Age             int       `json:"age"`
}

func (u *User) TableName() string {
	return "users"
}
