package download

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
	group := router.Group("/api/download")
	group.GET("/file", module.CodecWrap(DownloadFile,
		codec.MakeCodec(codec.DefaultStreamCodec, codec.DefaultURLCodec), &models.DownloadFileRequest{}, module.SecretAuth))
	group.GET("/block", module.CodecWrap(DownloadBlock,
		codec.MakeCodec(codec.DefaultStreamCodec, codec.DefaultURLCodec), &models.DownloadBlockRequest{}, module.SecretAuth))
}
