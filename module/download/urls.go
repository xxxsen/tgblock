package download

import (
	"tgblock/module"

	"github.com/gin-gonic/gin"
)

func init() {
	module.Regist(InitModule)
}

func InitModule(router *gin.Engine) {
	group := router.Group("/api/download")
	group.GET("/file", module.StreamWrap(DownloadFile, module.DefaultURLCodec, &DownloadFileRequest{}))
}
