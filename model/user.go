package model

import "gorm.io/gorm"

type User struct {
	Name      string  `gorm:"varchar(15);unique;not null"`
	Password  string  `gorm:"varchar(20);not null"`
	Gender    string  `gorm:"default:'man'"`
	Email     string  `gorm:"type:varchar(20);unique index"`
	Age       uint    `gorm:"default:0"`
	Birthday  int64   `gorm:"default:20000101"`
	Signature string  `gorm:"type:varchar(100)"`
	HeadPhoto string  `gorm:"type:varchar(100)"`
	Friends   []User  `gorm:"many2many:user_friends"`
	Groups    []Group `gorm:"many2many:user_group"`
	gorm.Model
}
