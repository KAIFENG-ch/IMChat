package model

import "gorm.io/gorm"

type Group struct {
	Name  string `gorm:"not null"`
	Users []User `gorm:"many2many:user_group"`
	gorm.Model
}
