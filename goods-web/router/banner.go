package router

import (
	"github.com/gin-gonic/gin"
	"hcshop-api/goods-web/api/banner"
	"hcshop-api/goods-web/middlewares"
)

func InitBannerRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("banner")
	{
		GoodsRouter.GET("", banner.List)
		GoodsRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), banner.New)
		GoodsRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), banner.Delete)
		GoodsRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), banner.Update)

	}
}
