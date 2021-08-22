package global

import (
	ut "github.com/go-playground/universal-translator"
	"hcshop-api/user_web/config"
	"hcshop-api/user_web/proto"
)

var (
	Trans         ut.Translator
	ServerConfig  = &config.ServerConfig{}
	NacosConfig   = &config.NacosConfig{}
	UserSrvClient proto.UserClient
)
