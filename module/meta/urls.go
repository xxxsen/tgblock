package meta

import (
	codec "github.com/xxxsen/tgblock/coder/server"

	"github.com/xxxsen/tgblock/module"
	"github.com/xxxsen/tgblock/module/models"

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
