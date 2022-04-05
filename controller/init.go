package controller

import (
	"IMChat/middleware"
	"IMChat/service"
	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	r := gin.Default()
	//r.Use(gin.Recovery(), gin.Logger())
	v1 := r.Group("/")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(200, "SUCCESS")
		})
		//v1.GET("ws", service.WsHandler)
		v1.POST("/register", Register)
		v1.POST("/login", Login)
		loginRequired := v1.Group("/user")
		loginRequired.Use(middleware.JWT())
		{
			loginRequired.GET("/ws", service.WsHandler)
			loginRequired.PUT("/update", Update)
			loginRequired.POST("/group", CreateGroup)
		}
	}
	return r
}
