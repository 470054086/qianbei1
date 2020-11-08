package routers

import (
	"github.com/gin-gonic/gin"
	"qianbei.com/app/controller"
)

func Category(r *gin.RouterGroup) {
	c := &controller.Category{}
	// 账本相关的路由
	g := r.Group("/category")
	g.Use()
	{
		g.POST("/index", c.Index)   //修改
		g.POST("/create", c.Create) //添加
		g.POST("/update", c.Update) //修改名称
		g.POST("/del", c.Del)       //删除
	}
}
