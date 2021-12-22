package meta

import (
	codec "tgblock/coder/server"
	"tgblock/module"
	"tgblock/module/models"

	"github.com/gin-gonic/gin"
)

func init() {
	module.Regist(InitModule)
}

func InitModule(router *gin.Engine) {
	group := router.Group("/api/meta")
	group.GET("/getfileinfo", module.CodecWrap(GetFileInfo,
		codec.MakeCodec(codec.DefaultJsonCodec, codec.DefaultURLCodec), &models.GetFileInfoRequest{}, module.SecretAuth))
}
