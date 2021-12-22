package download

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
	group := router.Group("/api/download")
	group.GET("/file", module.CodecWrap(DownloadFile,
		codec.MakeCodec(codec.DefaultStreamCodec, codec.DefaultURLCodec), &models.DownloadFileRequest{}, module.TokenAuth))
	group.GET("/block", module.CodecWrap(DownloadBlock,
		codec.MakeCodec(codec.DefaultStreamCodec, codec.DefaultURLCodec), &models.DownloadBlockRequest{}, module.TokenAuth))
}
