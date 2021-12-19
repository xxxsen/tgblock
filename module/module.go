package module

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"tgblock/module/constants"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/log"
)

var defaultCtx = &ServiceContext{}
var defaultEngine = gin.New()

type CGIHandler func(sctx *ServiceContext, gin *gin.Context, request interface{}) (int, interface{}, error)

func tryDecodeRequest(gctx *gin.Context, pType interface{}) (interface{}, error) {
	if pType == nil {
		return nil, nil

	}
	data, err := ioutil.ReadAll(gctx.Request.Body)
	if err != nil {
		return nil, err
	}
	dataType := reflect.TypeOf(pType).Elem()
	ptr := reflect.New(dataType)
	if err := json.Unmarshal(data, ptr.Interface()); err != nil {
		return nil, err
	}
	return ptr.Interface(), nil
}

func JsonWrap(handler CGIHandler, params interface{}) gin.HandlerFunc {
	return func(gctx *gin.Context) {
		paramsInstance, err := tryDecodeRequest(gctx, params)
		if err != nil {
			gctx.AbortWithStatusJSON(http.StatusBadRequest, GinErrResponse(constants.ErrUnknown, "decode request fail", err))
			return
		}
		code, data, err := handler(defaultCtx, gctx, paramsInstance)
		if code == http.StatusOK {
			gctx.JSON(code, GinResponse(data))
			return
		}
		log.Errorf("API:%s exec fail, code:%d, err:%v", gctx.Request.URL.Path, code, err)
		if e, ok := err.(*APIError); ok {
			gctx.AbortWithStatusJSON(code, GinErrResponse(e.Code, e.Errmsg, e.Err))
			return
		}
		gctx.AbortWithError(code, err)
	}
}

type ModuleInitFunc func(*gin.Engine)

func Regist(caller ModuleInitFunc) {
	caller(defaultEngine)
}

func Init(opts ...Option) error {
	for _, opt := range opts {
		opt(defaultCtx)
	}
	return nil
}

func Run(address string) error {
	for _, item := range defaultEngine.Routes() {
		log.Debugf("URI:%s, METHOD:%s", item.Path, item.Method)
	}
	return defaultEngine.Run(address)
}
