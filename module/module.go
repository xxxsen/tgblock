package module

import (
	"github.com/gin-gonic/gin"
	"github.com/xxxsen/log"
)

var defaultCtx = &ServiceContext{}
var defaultEngine = gin.New()

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
