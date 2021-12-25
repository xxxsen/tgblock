package share

import (
	"tgblock/module"
	"tgblock/module/models"

	codec "tgblock/coder/server"

	"github.com/gin-gonic/gin"
)

func init() {
	module.Regist(InitModule)
}

func InitModule(router *gin.Engine) {
	group := router.Group("/api/share")
	group.POST("/createshare", module.CodecWrap(CreateShare,
		codec.DefaultJsonCodec, &models.CreateShareRequest{}, module.SecretAuth))
	group.GET("/getshare", module.CodecWrap(GetShare,
		codec.MakeCodec(codec.DefaultStreamCodec, codec.DefaultURLCodec), &models.GetShareRequest{}, module.NoAuth))
}
