package routers

import (
	"github.com/gin-gonic/gin"
	"qianbei.com/app/controller"
)

func Book(r *gin.RouterGroup) {
	c := &controller.Book{}
	// 账本相关的路由
	g := r.Group("/book")
	g.Use()
	{
		g.POST("/index", c.Index)   //查询
		g.POST("/create", c.Create) // 创建
		g.POST("/update", c.Update) //修改
		g.POST("/join", c.JoinUser) // 加入
		g.POST("/del", c.Delete)    // 删除
	}
}
