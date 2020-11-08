package main

import (
	"qianbei.com/core" //自动加载init
	"qianbei.com/initialize"
)

func main() {
	defer core.Redis().Close()
	defer core.Db().Close()
	// 启动路由 Http服务
	initialize.InitHttpServer()
}
