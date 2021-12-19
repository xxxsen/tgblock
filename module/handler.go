package module

import (
	"io"
	"net/http"
	"tgblock/module/constants"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/log"
)

type CGIHandler func(sctx *ServiceContext, ctx *gin.Context, request interface{}) (int, interface{}, error)
type StreamHandler func(sctx *ServiceContext, ctx *gin.Context, request interface{}) (int, io.ReadCloser, error)

func CodecWrap(handler CGIHandler, codec Codec, params interface{}) gin.HandlerFunc {
	return func(gctx *gin.Context) {
		inst, err := codec.Decode(gctx, params)
		if err != nil {
			gctx.AbortWithStatusJSON(http.StatusBadRequest, GinErrResponse(constants.ErrUnknown, "decode request fail", err))
			return
		}
		code, data, err := handler(defaultCtx, gctx, inst)
		if code == http.StatusOK {
			gctx.JSON(code, GinResponse(data))
			return
		}
		log.Errorf("API:%s exec fail, code:%d, err:%v", gctx.Request.URL.Path, code, err)
		e := AsAPIError(err)
		gctx.AbortWithStatusJSON(code, GinErrResponse(e.Code, e.Errmsg, e.Err))
	}
}

func JsonWrap(handler CGIHandler, params interface{}) gin.HandlerFunc {
	return CodecWrap(handler, DefaultJsonCodec, params)
}

func URLWrap(handler CGIHandler, params interface{}) gin.HandlerFunc {
	return CodecWrap(handler, DefaultURLCodec, params)
}

func StreamWrap(handler StreamHandler, codec Codec, params interface{}) gin.HandlerFunc {
	return func(gctx *gin.Context) {
		inst, err := codec.Decode(gctx, params)
		if err != nil {
			gctx.AbortWithStatusJSON(http.StatusBadRequest, GinErrResponse(constants.ErrUnknown, "decode request fail", err))
			return
		}
		code, rc, err := handler(defaultCtx, gctx, inst)
		if code != http.StatusOK {
			e := AsAPIError(err)
			log.Errorf("API:%s exec fail, code:%d, err:%v", gctx.Request.URL.Path, code, err)
			gctx.AbortWithStatusJSON(code, GinErrResponse(e.Code, e.Errmsg, e.Err))
			return
		}
		defer rc.Close()
		_, err = io.Copy(gctx.Writer, rc)
		if err != nil {
			log.Errorf("copy stream fail, url:%s, code:%d, err:%v", gctx.Request.URL.Path, code, err)
		}
	}
}
