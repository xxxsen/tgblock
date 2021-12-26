package module

import (
	"net/http"
	"reflect"

	"github.com/xxxsen/tgblock/coder/errs"
	"github.com/xxxsen/tgblock/coder/frame"

	coder "github.com/xxxsen/tgblock/coder/server"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/log"
)

type CGIHandler func(sctx *ServiceContext, ctx *gin.Context, request interface{}) (int, interface{}, error)

func newInst(params interface{}) interface{} {
	if params == nil {
		return nil
	}
	t := reflect.TypeOf(params)
	v := reflect.New(t.Elem())
	return v.Interface()
}

func doAuth(auth CommonAuth, sctx *ServiceContext, gctx *gin.Context) bool {
	ok, err := auth.Auth(sctx, gctx)
	if err != nil {
		log.Errorf("url:%s, auth fail, err:%v", gctx.Request.URL.Path, err)
		return false
	}
	return ok
}

func CodecWrap(handler CGIHandler, codec coder.Codec, params interface{}, auth CommonAuth) gin.HandlerFunc {
	return func(gctx *gin.Context) {
		if !doAuth(auth, defaultCtx, gctx) {
			gctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		inst := newInst(params)
		if inst != nil {
			if err := codec.Decode(gctx, inst); err != nil {
				log.Errorf("decode request fail, url:%s, err:%v", gctx.Request.URL.Path, err)
				e := errs.AsAPIError(err)
				gctx.JSON(http.StatusBadRequest, frame.MakeErrJsonFrame(e.Code, e.Error()))
				return
			}
		}
		code, data, errHandle := handler(defaultCtx, gctx, inst)
		errEncode := codec.Encode(gctx, code, data, errHandle)
		if errHandle != nil || errEncode != nil || code != http.StatusOK {
			log.Errorf("handle request fail, url:%s, code:%d, err.handle:%v, err.encode:%v", gctx.Request.URL.Path, code, errHandle, errEncode)
			return
		}
	}
}
