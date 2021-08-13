package main

import (
	"fmt"
	"go.uber.org/zap"
	"hcshop-api/user_web/initialize"
)

func main() {
	port := 8021
	// 日志初始化
	initialize.InitLogger()
	// 路由初始化
	Router := initialize.Routers()
	zap.S().Infof("启动服务,端口%d", port)
	if err := Router.Run(fmt.Sprintf(":%d", port)); err != nil {
		zap.S().Panic("启动失败:", err.Error())
	}
}
