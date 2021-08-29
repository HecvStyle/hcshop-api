package router

import (
	"github.com/gin-gonic/gin"
	"hcshop-api/goods-web/api/goods"
)

func InitGoodsRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("goods")
	{
		UserRouter.GET("/list", goods.List)
	}
}
