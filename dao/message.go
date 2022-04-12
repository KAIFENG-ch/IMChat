package dao

import (
	"IMChat/model"
	"strconv"
	"time"
)

func InsertMsg(uid int, toUid int, content string, expire int64, status bool) {
	tx := model.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	insertMsg := &model.Message{
		Status:   status,
		UserID:   uid,
		ToUserID: toUid,
		Content:  content,
		EndAt:    time.Now().Unix() + expire,
	}
	model.DB.Create(&insertMsg)
}

func ReadMessage(uid string, toUid string) []model.Message {
	var msg []model.Message
	UID, _ := strconv.Atoi(uid)
	toUID, _ := strconv.Atoi(toUid)
	model.DB.Model(&model.Message{}).Order("created_at desc").Limit(5).
		Where("user_id = ? and to_user_id = ?", UID, toUID).
		Or("user_id = ? and to_user_id = ?", toUID, UID).Find(&msg)
	return msg
}
