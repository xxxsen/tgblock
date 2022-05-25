package upload

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
	group := router.Group("/api/upload")
	group.POST("/post", module.CodecWrap(PostUpload, codec.MakeEncoder(codec.DefaultJsonCodec), nil))
	group.POST("/block/begin", module.CodecWrap(BlockUploadBegin, codec.DefaultJsonCodec, &models.BlockUploadBeginRequest{}))
	group.POST("/block/part", module.CodecWrap(BlockUploadPart, codec.MakeEncoder(codec.DefaultJsonCodec), nil))
	group.POST("/block/end", module.CodecWrap(BlockUploadEnd, codec.DefaultJsonCodec, &models.BlockUploadEndRequest{}))
}
