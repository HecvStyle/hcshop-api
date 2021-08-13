package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"hcshop-api/user_web/global/response"
	"hcshop-api/user_web/proto"
	"net/http"
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

func GetUserList(ctx *gin.Context) {
	ip := "127.0.0.1"
	port := 50051
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务失败】", "msg", err.Error())
	}
	userClient := proto.NewUserClient(conn)

	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    1,
		PSize: 10,
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
