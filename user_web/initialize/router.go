package initialize

import (
	"github.com/gin-gonic/gin"
	"hcshop-api/user_web/middlewares"
	"hcshop-api/user_web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	// 为所有请求配置跨域处理
	Router.Use(middlewares.Cors())
	ApiGroup := Router.Group("/u/v1")
	router.InitUserRouter(ApiGroup)
	router.InitBaseRouter(ApiGroup)
	return Router
}
