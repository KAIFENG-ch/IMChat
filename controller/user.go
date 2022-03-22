package controller

import (
	"IMChat/service"
	"github.com/gin-gonic/gin"
)

func Register(ctx *gin.Context) {
	var register service.UserRegister
	err := ctx.ShouldBind(&register)
	if err != nil {
		return
	}
	res := register.Register()
	ctx.JSON(200, res)
}

func Login(c *gin.Context) {
	var login service.UserRegister
	err := c.ShouldBind(&login)
	if err != nil {
		return
	}
	res := login.Login()
	c.JSON(200, res)
}
