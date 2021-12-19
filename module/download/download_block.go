package download

import (
	"io"
	"net/http"
	"tgblock/module"
	"tgblock/module/constants"
	"tgblock/processor"
	"tgblock/shortten"

	"github.com/gin-gonic/gin"
)

func DownloadBlock(sctx *module.ServiceContext, ctx *gin.Context, params interface{}) (int, io.ReadCloser, error) {
	req := params.(*DownloadBlockRequest)
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
	if int(req.BlockIndex) >= len(meta.FileList) {
		return http.StatusBadRequest, nil, module.NewAPIError(constants.ErrParams, "index out of range")
	}
	return http.StatusOK, newPartReader(sctx, meta.FileList[req.BlockIndex].FileId), nil
}
