package upload

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
	group := router.Group("/api/upload")
	group.POST("/post", module.CodecWrap(PostUpload, codec.MakeEncoder(codec.DefaultJsonCodec), nil, module.SecretAuth))
	group.POST("/block/begin", module.CodecWrap(BlockUploadBegin, codec.DefaultJsonCodec, &models.BlockUploadBeginRequest{}, module.SecretAuth))
	group.POST("/block/part", module.CodecWrap(BlockUploadPart, codec.MakeEncoder(codec.DefaultJsonCodec), nil, module.SecretAuth))
	group.POST("/block/end", module.CodecWrap(BlockUploadEnd, codec.DefaultJsonCodec, &models.BlockUploadEndRequest{}, module.SecretAuth))
}
