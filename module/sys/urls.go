package sys

import (
	"github.com/xxxsen/tgblock/module"
	"github.com/xxxsen/tgblock/module/models"

	codec "github.com/xxxsen/tgblock/coder/server"

	"github.com/gin-gonic/gin"
)

func init() {
	module.Regist(InitModule)
}

func InitModule(router *gin.Engine) {
	group := router.Group("/api/sys")
	group.GET("/getsysinfo", module.CodecWrap(GetSysInfo,
		codec.MakeCodec(codec.DefaultJsonCodec, codec.DefaultURLCodec), &models.GetSysInfoRequest{}, module.SecretAuth))

}
