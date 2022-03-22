package controller

import (
	"IMChat/service"
	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery(), gin.Logger())
	v1 := r.Group("/")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(200, "SUCCESS")
		})
		v1.GET("ws", service.WsHandler)
		v1.POST("/register", Register)
	}
	return r
}
