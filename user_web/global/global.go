package global

import (
	ut "github.com/go-playground/universal-translator"
	"hcshop-api/user_web/config"
)

var (
	Trans        ut.Translator
	ServerConfig = &config.ServerConfig{}
)
