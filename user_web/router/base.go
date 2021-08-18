package router

import (
	"github.com/gin-gonic/gin"
	"hcshop-api/user_web/api"
)

func InitBaseRouter(Router *gin.RouterGroup) {
	BaseRouter := Router.Group("base")
	{
		BaseRouter.GET("captcha", api.GetCaptcha)
		BaseRouter.POST("sendSms", api.SendSms)
	}
}
