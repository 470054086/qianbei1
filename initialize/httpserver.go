package initialize

// 路由初始化
import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"qianbei.com/app/middleware"
	"qianbei.com/core"
	"qianbei.com/routers"
)

// 使用gin的路由
func InitHttpServer() {
	e := gin.New()
	// 定义全局中间件
	e.Use(middleware.Error())
	// 添加路由文件
	addRouter(e.Group(""))
	port := core.Config().App.Port
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), e)
	core.QLog().Info("http port %s is running", port)
	if err != nil {
		panic(fmt.Sprintf("http start is erros %v", err))
	}
}

// 添加路由文件
func addRouter(group *gin.RouterGroup) {
	// 全局错误中间件
	routers.Pay(group)
	routers.Book(group)
	routers.Category(group)
}
