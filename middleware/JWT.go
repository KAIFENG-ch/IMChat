package middleware

import (
	"IMChat/utils"
	"github.com/gin-gonic/gin"
	"time"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		token := c.GetHeader("Authorization")
		if token == "" {
			code = 404
		} else {
			claims, err := utils.ParseToken(token)
			if err != nil {
				code = 30001
			} else if claims.ExpiresAt < time.Now().Unix() {
				code = 30002
			}
			code = 200
		}
		if code != 200 {
			c.JSON(200, gin.H{
				"status": 200,
				"msg":    "token error",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
