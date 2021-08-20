package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"hcshop-api/user_web/global"
	"hcshop-api/user_web/proto"
)

func InitSrvConn() {
	consulInfo := global.ServerConfig.ConsulInfo
	conn, err := grpc.Dial(
		// 记得添加 consul:// 协议，别搞错了 踩坑+1
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("用户服务注册连接失败")
	}
	// 这里涉及到了多个链接，都只用了一个gorutine,考虑使用连接池才可以
	global.UserSrvClient = proto.NewUserClient(conn)

}

func InitSrvConn2() {
	cfg := api.DefaultConfig()
	consulInfo := global.ServerConfig.ConsulInfo
	cfg.Address = fmt.Sprintf("%s:%d", consulInfo.Host, consulInfo.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	userSrvHost := ""
	userSrvPort := 0
	// 这里注意过滤的语法， ==  右边是需要 双引号包裹的。 坑+1
	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, global.ServerConfig.UserSrvInfo.Name))
	//data, err := client.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%s\"", global.ServerConfig.UserSrvInfo.Name))
	// 这里缺了双引号会报错的，直接panic
	//data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == %s`, global.ServerConfig.UserSrvInfo.Name))

	if err != nil {
		panic(err)
	}
	for _, value := range data {
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}

	if userSrvHost == "" {
		zap.S().Fatal("【InitSrvConn】 获取用户服务失败")
		return
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务失败】", "msg", err.Error())
	}

	// 这里涉及到了多个链接，都只用了一个gorutine,考虑使用连接池才可以
	global.UserSrvClient = proto.NewUserClient(conn)

}
