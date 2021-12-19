package meta

import (
	"tgblock/module"

	"github.com/gin-gonic/gin"
)

func init() {
	module.Regist(InitModule)
}

func InitModule(router *gin.Engine) {
	group := router.Group("/api/meta")
	group.GET("/getfileinfo", module.URLWrap(GetFileInfo, &GetFileInfoRequest{}))
}
