package dao

import (
	"IMChat/model"
	"context"
	"time"
)

type Trainer struct {
	Content string `json:"content"`
	StartAt int64  `json:"start_at"`
	EndAt   int64  `json:"end_at"`
	Read    uint   `json:"read"`
}

func InsertMsg(database string, id string, content string, read uint, expire int64) (err error) {
	collection := model.MDB.Database(database).Collection(id)
	comment := &Trainer{
		Content: content,
		StartAt: time.Now().Unix(),
		EndAt:   time.Now().Unix() + expire,
		Read:    read,
	}
	_, err = collection.InsertOne(context.TODO(), comment)
	return
}

//var results []Trainer
//
//func Read(database string, id string, time int64, content string, read uint, expire int64) (err error) {
//	idCollect := model.MDB.Database(database).Collection(id)
//	filter := bson.M{"startAt": bson.M{"$lt": time}}
//	sendIdCursor,err := idCollect.Find(context.TODO(), filter, options.Find().SetSort(bson.D{{"StartAt": -1}}))
//}
