package controller

import (
	"IMChat/service"
	"IMChat/utils"
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
		c.JSON(400, err)
	}
	res := login.Login()
	c.JSON(200, res)
}

func Update(c *gin.Context) {
	var updates service.UserUpdate
	claims, _ := utils.ParseToken(c.GetHeader("Authorization"))
	err := c.ShouldBind(&updates)
	if err != nil {
		c.JSON(400, err)
	}
	res := updates.Update(claims.Id)
	c.JSON(200, res)
}

func CreateGroup(c *gin.Context) {
	var group service.GroupRegister
	claims, _ := utils.ParseToken(c.GetHeader("Authorization"))
	err := c.ShouldBind(&group)
	if err != nil {
		c.JSON(400, err)
	}
	res := group.CreateGroup(claims.Id)
	c.JSON(200, res)
}
