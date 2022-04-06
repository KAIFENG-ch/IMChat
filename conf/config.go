package conf

import (
	"IMChat/model"
	"fmt"
	"github.com/spf13/viper"
)

type SqlDb struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
}

type RedisDb struct {
	Host string
	Port string
}

type MongoDb struct {
	Host     string
	Port     string
	Database string
}

type Dbs struct {
	Sql   SqlDb
	Redis RedisDb
	Mongo MongoDb
}

func InitConfig() {
	viper.SetConfigType("yaml")
	viper.SetConfigFile("./conf/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return
	}
	var myConfig Dbs
	err = viper.Unmarshal(&myConfig)
	if err != nil {
		return
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?"+
		"charset=utf8mb4&parseTime=True&loc=Local", myConfig.Sql.User,
		myConfig.Sql.Password, myConfig.Sql.Host, myConfig.Sql.Port, myConfig.Sql.Database)
	model.Database(dsn)
	err = model.RedisDB(myConfig.Redis.Host, myConfig.Redis.Port)
	if err != nil {
		return
	}
}
