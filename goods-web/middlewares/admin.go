package middlewares

import (
	"github.com/gin-gonic/gin"
	"hcshop-api/goods-web/models"
	"net/http"
)

func IsAdminAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		claims, _ := context.Get("claims")
		currentUser := claims.(*models.CustomClaims)
		if currentUser.AuthorityId != 2 {
			context.JSON(http.StatusForbidden, gin.H{
				"msg": "没有权限",
			})
			context.Abort()
			return
		}
		context.Next()
	}
}
