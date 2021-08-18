package api

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"hcshop-api/user_web/forms"
	"hcshop-api/user_web/global"
	"hcshop-api/user_web/global/response"
	"hcshop-api/user_web/middlewares"
	"hcshop-api/user_web/models"
	"hcshop-api/user_web/proto"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func HandlerGrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": e.Code(),
				})

			}
		}
	}
}

func HandleValidatorErr(ctx *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
}
func GetUserList(ctx *gin.Context) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务失败】", "msg", err.Error())
	}
	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	zap.S().Infof("访问用：%d", currentUser.ID)

	userClient := proto.NewUserClient(conn)
	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)
	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询 【用户列表】失败")
		HandlerGrpcErrorToHttp(err, ctx)
		return
	}
	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		data := response.UserResponse{}
		data.Id = value.Id
		data.Nickname = value.Nickname
		//data.Birthday = time.Time(time.Unix(int64(value.Birthday), 0)).Format("2006-01-02")
		// 这个是时间格式化的处理
		data.Birthday = response.JsonTime(time.Unix(int64(value.Birthday), 0))
		data.Gender = value.Gender
		data.Mobile = value.Mobile
		result = append(result, data)
	}
	ctx.JSON(http.StatusOK, result)
}

func removeTopStruct(fields map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fields {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

func PassWordLogin(ctx *gin.Context) {
	passwordLoginForm := forms.PasswordLoginForm{}
	if err := ctx.ShouldBind(&passwordLoginForm); err != nil {
		HandleValidatorErr(ctx, err)
		return
	}
	if !store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, true) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"captcha": "验证码错误",
		})
		return
	}
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务失败】", "msg", err.Error())
		return
	}
	userClient := proto.NewUserClient(conn)

	if resp, err := userClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	}); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusBadRequest, gin.H{
					"msg": "用户不存在",
				})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": "登录失败",
				})
			}
			return
		}
	} else {
		if pwd_rsp, err := userClient.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
			Password:          passwordLoginForm.Password,
			EncryptedPassword: resp.Password,
		}); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": "登录失败",
			})
		} else {
			if pwd_rsp.Success {
				// 这里要生成token
				j := middlewares.NewJWT()
				claims := models.CustomClaims{
					ID:          uint(resp.Id),
					NickName:    resp.Nickname,
					AuthorityId: uint(resp.Role),
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Unix(), // 签名生效时间
						ExpiresAt: time.Now().Unix() + 60*60*24*30,
						Issuer:    "hecv",
					},
				}

				token, err := j.CreateToken(claims)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"msg": "生成token失败",
					})
					return
				}
				ctx.JSON(http.StatusOK, gin.H{
					"token":      token,
					"id":         resp.Id,
					"nickname":   resp.Nickname,
					"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
				})
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"msg": "登录失败",
				})
			}
		}
	}
}
