package sys

import (
	"net/http"
	"tgblock/module"
	"tgblock/module/models"

	"github.com/gin-gonic/gin"
)

func GetSysInfo(sctx *module.ServiceContext, ctx *gin.Context, params interface{}) (int, interface{}, error) {
	rsp := &models.GetSysInfoResponse{
		MaxFileSize: sctx.MaxFileSize,
		BlockSize:   sctx.BlockSize,
	}
	return http.StatusOK, rsp, nil
}
