package routers

import (
	"github.com/gin-gonic/gin"
	"qianbei.com/app/controller"
)

func Pay(r *gin.RouterGroup) {
	c := &controller.Pay{}
	g := r.Group("/pay")
	g.Use()
	{
		g.POST("/record", c.AddRecord)
		g.POST("/index",c.Record)
	}
}
