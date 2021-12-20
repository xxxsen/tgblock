package sys

import (
	"net/http"
	"tgblock/module"

	"github.com/gin-gonic/gin"
)

func GetSysInfo(sctx *module.ServiceContext, ctx *gin.Context, params interface{}) (int, interface{}, error) {
	rsp := &GetSysInfoResponse{
		MaxFileSize: sctx.MaxFileSize,
		BlockSize:   sctx.BlockSize,
	}
	return http.StatusOK, rsp, nil
}
