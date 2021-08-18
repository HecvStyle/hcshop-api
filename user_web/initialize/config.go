package initialize

import (
	"fmt"
	"hcshop-api/user_web/global"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	// 通过主机的环境变量，来区分线上环境和本地环境，
	// ！！！需要自己配置本地环境变量 HCSHOP_DEBUG ！！！
	debug := GetEnvInfo("HCSHOP_DEBUG")
	configFilePrefix := "config"
	configName := fmt.Sprintf("%s-pro.yaml", configFilePrefix)
	if debug {
		configName = fmt.Sprintf("%s-debug.yaml", configFilePrefix)
	}
	v := viper.New()
	v.SetConfigFile(configName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("配置信息：%v", global.ServerConfig)
	// 这里是监控配置文件改变
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infof("配置文件发生:%s", e.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.ServerConfig)
		zap.S().Infof("配置信息：%v", global.ServerConfig)
	})

}
