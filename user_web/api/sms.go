package api

import (
	"context"
	"fmt"
	//"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	//"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"hcshop-api/user_web/forms"
	"hcshop-api/user_web/global"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// GenerateSmsCode 短信验证码生成
func GenerateSmsCode(length int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())
	var sb strings.Builder
	for i := 0; i < length; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

// SendSms 发送验证码
func SendSms(ctx *gin.Context) {
	sendSmsForm := forms.SendSmsForm{}
	if err := ctx.ShouldBind(&sendSmsForm); err != nil {
		HandleValidatorErr(ctx, err)
		return
	}

	//client, err := dysmsapi.NewClientWithAccessKey("cn-beijing", global.ServerConfig.AliSmsInfo.ApiKey, global.ServerConfig.AliSmsInfo.ApiSecret)
	//if err != nil {
	//	panic(err)
	//}
	smsCode := GenerateSmsCode(6)
	//request := requests.NewCommonRequest()
	//request.Method = "POST"
	//request.Scheme = "https" // https | http
	//request.Domain = "dysmsapi.aliyuncs.com"
	//request.Version = "2017-05-25"
	//request.ApiName = "SendSms"
	//request.QueryParams["RegionId"] = "cn-beijing"
	//request.QueryParams["PhoneNumbers"] = sendSmsForm.Mobile
	//request.QueryParams["SignName"] = "go测试"
	//request.QueryParams["TemplateCode"] = "SMS_181850725"
	//request.QueryParams["TemplateParam"] = "{\"code\":" + smsCode + "}" //短信模板中的验证码内容 自己生成   之前试过直接返回，但是失败，加上code成功。
	//response, err := client.ProcessCommonRequest(request)
	//fmt.Print(client.DoAction(request, response))
	//if err != nil {
	//	fmt.Print(err.Error())
	//}

	// TODO: 这里无法提供私人的阿里云AK，所以生产验证码直接返回测试
	// ⚠️⚠️⚠️⚠️  这是不安全的做法，生产环境上绝对的严重漏洞  ⚠️⚠️⚠️⚠️⚠️

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
		Password: global.ServerConfig.RedisInfo.Password,
		DB:       0,
	})

	err := rdb.Set(context.Background(), sendSmsForm.Mobile, smsCode, time.Duration(global.ServerConfig.AliSmsInfo.Expire)*time.Minute).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生产验证码出错",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "发送成功",
		"code": smsCode,
	})
	return
}
