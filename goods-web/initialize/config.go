package initialize

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"hcshop-api/goods-web/config"
	"hcshop-api/goods-web/global"

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
	if err := v.Unmarshal(global.NacosConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("配置信息：%v", global.NacosConfig)
	// 这里是监控配置文件改变
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infof("配置文件发生:%s", e.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.NacosConfig)
		zap.S().Infof("配置信息：%v", global.NacosConfig)
	})

	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   global.NacosConfig.Port,
		},
	}
	clientConfig := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		//TODO:  因为我这里gland用的工作环境是goods-web 目录，考虑后期nacos的缓存文件应该放到根目录下，所以之类指定了工作目录的上一程，这个可能在正式环境需要改变，所以考虑写进配置文件中来
		LogDir:     "../tmp/nacos/log",
		CacheDir:   "../tmp/nacos/cache",
		RotateTime: "1h",
		MaxAge:     3,
		LogLevel:   "debug",
	}

	// 创建动态配置客户端的另一种方式 (推荐)
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)

	if err != nil {
		panic(err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
	})
	if err != nil {
		panic(err)
	}

	serverConfig := config.ServerConfig{}
	err = json.Unmarshal([]byte(content), &serverConfig)
	if err != nil {
		zap.S().Fatalf("读取nacos配置失败: %s", err.Error())
	}
	global.ServerConfig = &serverConfig

}
