package sys

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
	group := router.Group("/api/sys")
	group.GET("/getsysinfo", module.CodecWrap(GetSysInfo, codec.MakeCodec(codec.DefaultJsonCodec, codec.DefaultURLCodec), &models.GetSysInfoRequest{}))

}
