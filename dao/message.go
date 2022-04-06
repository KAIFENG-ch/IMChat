package dao

import (
	"IMChat/model"
	"time"
)

func InsertMsg(uid int, content string, expire int64, status bool) {
	insertMsg := &model.Message{
		Status:  status,
		UserID:  uid,
		Content: content,
		EndAt:   time.Now().Unix() + expire,
	}
	model.DB.Create(&insertMsg)
}
