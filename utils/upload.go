package utils

import (
	"bytes"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/viper"
	"log"
	"path/filepath"
	"time"
)

func UploadToCloud(object string, fileByte []byte) (url string, err error) {
	var endPoint = viper.GetString("OSS.EndPoint")
	var accessKey = viper.GetString("OSS.AccessKey")
	var SecretKey = viper.GetString("OSS.SecretKey")
	var bucket = viper.GetString("OSS.Bucket")
	var region = viper.GetString("OSS.Region")

	client, err := oss.New(endPoint, accessKey, SecretKey)
	if err != nil {
		log.Println(err)
		return "", err
	}
	myBucket, err := client.Bucket(bucket)
	if err != nil {
		log.Println(err)
		return "", err
	}
	foldTime := time.Now().Format("2006-01-02")
	yunFilePath := filepath.Join("uploads", foldTime) + "/" + object
	err = myBucket.PutObject(object, bytes.NewReader(fileByte))
	if err != nil {
		return url, err
	}
	return region + "/" + yunFilePath + "/" + foldTime, nil
}
