package dao

import "IMChat/model"

func AddSet(id string, element string) {
	model.RDB.SAdd(id, element)
}

func IsMember(name string, key string) bool {
	return model.RDB.SIsMember(name, key).Val()
}
