package model

import "gorm.io/gorm"

type Message struct {
	UserID   int    `gorm:"not null"`
	ToUserID int    `gorm:"not null"`
	Content  string `gorm:"type:varchar(1000);not null"`
	EndAt    int64
	Status   bool `gorm:"default:false"`
	gorm.Model
}

type GroupMessage struct {
	GroupID  int    `gorm:"not null"`
	SenderID int    `gorm:"not null"`
	Content  string `gorm:"type:varchar(1000);not null"`
	EndAt    int64
	Status   bool `gorm:"default:false"`
}
