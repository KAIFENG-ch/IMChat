package model

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var DB *gorm.DB

func Database(dsn string) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}
	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxIdleTime(time.Hour)
	DB = db
	err = DB.Set("gorm:table_option", "ENGINE=InnoDB").
		AutoMigrate(&User{})
	if err != nil {
		return
	}
}

var RDB *redis.Client

func RedisDB(addr string, password string) (err error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	RDB = rdb
	_, err = RDB.Ping().Result()
	if err != nil {
		return
	}
	return nil
}

var MDB *mongo.Client

func MongoDB(dsn string) (e error) {
	ctx := context.TODO()
	opt := new(options.ClientOptions)
	opt = opt.SetMaxPoolSize(100)
	du, _ := time.ParseDuration("5000")
	opt = opt.SetConnectTimeout(du)
	MongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn), opt)
	MDB = MongoClient
	if err != nil {
		e = err
		return
	}
	e = MongoClient.Ping(ctx, nil)
	defer func(MongoClient *mongo.Client, ctx context.Context) {
		err := MongoClient.Disconnect(ctx)
		if err != nil {
			e = err
			return
		}
	}(MongoClient, ctx)
	databases, err := MongoClient.ListDatabases(ctx, bson.M{})
	if err != nil {
		e = err
		return
	}
	fmt.Println(databases)
	return nil
}
