package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"hcshop-api/user_web/global"
	"hcshop-api/user_web/initialize"
	"hcshop-api/user_web/utils"
	vl "hcshop-api/user_web/validator"
	"os"
	"os/signal"
	"syscall"
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

	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	//生成对应的检查对象
	check := &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/health", global.ServerConfig.Host, global.ServerConfig.Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	//生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServerConfig.Name
	registration.ID, _ = uuid.GenerateUUID()
	registration.Port = global.ServerConfig.Port
	registration.Tags = global.ServerConfig.Tags
	registration.Address = global.ServerConfig.Host
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}

	zap.S().Infof("用户web层启动服务,端口%d", global.ServerConfig.Port)
	go func() {
		if err = Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.S().Panic("启动失败:", err.Error())
		}
	}()
	zap.S().Info("用户web层服务启动完成")

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = client.Agent().ServiceDeregister(registration.ID); err != nil {
		zap.S().Info("用户web服务注销失败")
		panic(err)
	}
	zap.S().Info("注销成功")

}
