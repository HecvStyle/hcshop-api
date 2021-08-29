package global

import (
	ut "github.com/go-playground/universal-translator"
	"hcshop-api/goods-web/config"
	"hcshop-api/goods-web/proto"
)

var (
	Trans          ut.Translator
	ServerConfig   = &config.ServerConfig{}
	NacosConfig    = &config.NacosConfig{}
	GoodsSrvClient proto.GoodsClient
)
