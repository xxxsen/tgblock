package download

import (
	"fmt"
	"io"
	"net/http"
	"tgblock/module"
	"tgblock/module/constants"
	"tgblock/processor"
	"tgblock/shortten"

	"github.com/gin-gonic/gin"
)

func DownloadFile(sctx *module.ServiceContext, ctx *gin.Context, params interface{}) (int, io.ReadCloser, error) {
	req := params.(*DownloadFileRequest)
	if len(req.FileId) == 0 {
		return http.StatusBadRequest, nil, module.NewAPIError(constants.ErrParams, "invalid fileid")
	}
	fileid, err := shortten.Decode(ctx, req.FileId)
	if err != nil {
		return http.StatusInternalServerError, nil, module.WrapError(constants.ErrUnMarshal, "decode fileid fail", err)
	}
	proc := processor.NewFileProcessor(sctx.Bot)
	meta, err := proc.GetFileMeta(ctx, fileid)
	if err != nil {
		return http.StatusInternalServerError, nil, module.WrapError(constants.ErrIO, "read file meta fail", err)
	}
	//write download info
	ctx.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", meta.FileSize))
	ctx.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", meta.Name))
	return http.StatusOK, newMultiBlockReader(sctx, meta), nil
}
