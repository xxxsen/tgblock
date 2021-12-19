package upload

import (
	"tgblock/module"

	"github.com/gin-gonic/gin"
)

func init() {
	module.Regist(InitModule)
}

func InitModule(router *gin.Engine) {
	group := router.Group("/api/upload")
	group.POST("/post", module.JsonWrap(PostUpload, nil))
	group.POST("/block/begin", module.JsonWrap(BlockUploadBegin, &BlockUploadBeginRequest{}))
	group.POST("/block/part", module.JsonWrap(BlockUploadPart, nil))
	group.POST("/block/end", module.JsonWrap(BlockUploadEnd, &BlockUploadEndRequest{}))
}
