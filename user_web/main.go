package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"hcshop-api/user_web/global"
	"hcshop-api/user_web/initialize"
	"hcshop-api/user_web/utils"
	vl "hcshop-api/user_web/validator"
)

func main() {

	// 日志初始化
	initialize.InitLogger()

	initialize.InitConfig()

	// 用户服务初始化
	initialize.InitSrvConn()

	// 路由初始化
	Router := initialize.Routers()

	if err := initialize.InitTrans("zh"); err != nil {
		panic(err)
	}

	viper.AutomaticEnv()
	debug := viper.GetBool("HCSHOP_DEBUG")
	// 线上环境，动态生成端口
	if !debug {
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", vl.ValidateMobile)
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号码！", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	zap.S().Infof("启动服务,端口%d", global.ServerConfig.Port)
	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动失败:", err.Error())
	}
}
