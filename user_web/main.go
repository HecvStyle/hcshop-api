package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"hcshop-api/user_web/global"
	"hcshop-api/user_web/initialize"
	vl "hcshop-api/user_web/validator"
)

func main() {

	// 日志初始化
	initialize.InitLogger()

	initialize.InitConfig()

	// 路由初始化
	Router := initialize.Routers()

	if err := initialize.InitTrans("zh"); err != nil {
		panic(err)
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("mobile", vl.ValidateMobile)
	}

	zap.S().Infof("启动服务,端口%d", global.ServerConfig.Port)
	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动失败:", err.Error())
	}
}
