package initialize

import (
	"github.com/gin-gonic/gin"
	"hcshop-api/goods-web/middlewares"
	"hcshop-api/goods-web/router"
	"net/http"
)

func Routers() *gin.Engine {
	Router := gin.Default()

	Router.GET("/health", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"code":http.StatusOK,
			"success":true,
		})
	})


	// 为所有请求配置跨域处理
	Router.Use(middlewares.Cors())
	ApiGroup := Router.Group("/g/v1")
	router.InitGoodsRouter(ApiGroup)
	return Router
}
