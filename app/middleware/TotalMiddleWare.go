package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func Total() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("我是统计开始")
		c.Next()
		fmt.Println("我是统计结束")
	}
}
