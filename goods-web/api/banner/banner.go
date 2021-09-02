package banner

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"
	"hcshop-api/goods-web/api/apiutil"
	"hcshop-api/goods-web/forms"
	"hcshop-api/goods-web/global"
	"hcshop-api/goods-web/proto"
	"net/http"
	"strconv"
)

func List(ctx *gin.Context) {
	gr, err := global.GoodsSrvClient.BannerList(context.Background(), &emptypb.Empty{})
	if err != nil {
		apiutil.HandlerGrpcErrorToHttp(err, ctx)
		return
	}
	resp := make([]interface{}, 0)
	for _, value := range gr.Data {
		reMap := make(map[string]interface{})
		reMap["id"] = value.Id
		reMap["index"] = value.Index
		reMap["image"] = value.Image
		reMap["url"] = value.Url
		resp = append(resp, reMap)
	}
	ctx.JSON(http.StatusOK, resp)
}

func New(ctx *gin.Context) {
	bannerForm := forms.BannerForm{}
	if err := ctx.ShouldBindJSON(&bannerForm); err != nil {
		apiutil.HandleValidatorErr(ctx, err)
	}

	gr, err := global.GoodsSrvClient.CreateBanner(context.Background(), &proto.BannerRequest{
		Index: int32(bannerForm.Index),
		Url:   bannerForm.Url,
		Image: bannerForm.Image,
	})
	if err != nil {
		apiutil.HandleValidatorErr(ctx, err)
		return
	}

	result := make(map[string]interface{})
	result["id"] = gr.Id
	result["index"] = gr.Index
	result["url"] = gr.Url
	result["image"] = gr.Image

	ctx.JSON(http.StatusOK, result)
}

func Update(ctx *gin.Context) {
	bannerForm := forms.BannerForm{}
	if err := ctx.ShouldBindJSON(&bannerForm); err != nil {
		apiutil.HandleValidatorErr(ctx, err)
		return
	}
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	_, err = global.GoodsSrvClient.UpdateBanner(context.Background(), &proto.BannerRequest{
		Id:    int32(i),
		Index: int32(bannerForm.Index),
		Url:   bannerForm.Url,
	})
	if err != nil {
		apiutil.HandlerGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	_, err = global.GoodsSrvClient.DeleteBanner(context.Background(), &proto.BannerRequest{Id: int32(i)})
	if err != nil {
		apiutil.HandlerGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, "")
}
