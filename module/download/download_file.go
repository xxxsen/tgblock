package download

import (
	"net/http"
	"tgblock/coder/errs"
	codec "tgblock/coder/server"
	"tgblock/module"
	"tgblock/module/constants"
	"tgblock/module/models"
	"tgblock/processor"
	"tgblock/shortten"

	"github.com/gin-gonic/gin"
)

func DownloadFile(sctx *module.ServiceContext, ctx *gin.Context, params interface{}) (int, interface{}, error) {
	req := params.(*models.DownloadFileRequest)
	if len(req.FileId) == 0 {
		return http.StatusBadRequest, nil, errs.NewAPIError(constants.ErrParams, "invalid fileid")
	}
	fileid, err := shortten.Decode(ctx, req.FileId)
	if err != nil {
		return http.StatusInternalServerError, nil, errs.WrapError(constants.ErrUnMarshal, "decode fileid fail", err)
	}
	proc := processor.NewFileProcessor(sctx.Bot)
	meta, err := proc.GetFileMeta(ctx, fileid)
	if err != nil {
		return http.StatusInternalServerError, nil, errs.WrapError(constants.ErrIO, "read file meta fail", err)
	}
	//write download info
	output := &codec.StreamInfo{
		Stream: newMultiBlockReader(sctx, meta),
		Size:   meta.FileSize,
		Name:   meta.Name,
	}
	return http.StatusOK, output, nil
}
