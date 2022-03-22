package main

import (
	"IMChat/conf"
	"IMChat/controller"
	"IMChat/service"
	"github.com/gin-gonic/gin"
)

func main() {
	conf.InitConfig()
	gin.SetMode(gin.ReleaseMode)
	go service.Manager.Connect()
	r := controller.Routers()
	_ = r.Run(":8000")
}
