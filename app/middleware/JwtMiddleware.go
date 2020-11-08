package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func Jwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("我是jwt开始")
		c.Next()
		fmt.Println("我是jwt结束")
	}
}
