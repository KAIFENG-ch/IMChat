package logs

import (
	"fmt"
	"io"
	"os"
	"time"
)

func CheckFile(f string) bool {
	_, err := os.Stat(f)
	if err != nil {
		return true
	}
	return false
}

func MakeLog(logType string, logs string) {
	var file *os.File
	fileName := "./file/logs/" + time.Now().Format("20060102") + ".logs"
	if CheckFile(fileName) {
		_, err := os.Open(fileName)
		if err != nil {
			fmt.Println("文件已存在")
		} else {
			_, err = os.Create(fileName)
			if err != nil {
				fmt.Println("创建失败")
			}
		}
		_, err = io.WriteString(file, logType+time.Now().Format("20060102")+logs+"\n")
		if err != nil {
			fmt.Println(err)
		}
	}
}
