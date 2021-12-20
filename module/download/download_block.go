package download

import (
	"net/http"
	"tgblock/coder/errs"
	codec "tgblock/coder/server"
	"tgblock/module"
	"tgblock/module/constants"
	"tgblock/processor"
	"tgblock/shortten"

	"github.com/gin-gonic/gin"
)

func DownloadBlock(sctx *module.ServiceContext, ctx *gin.Context, params interface{}) (int, interface{}, error) {
	req := params.(*DownloadBlockRequest)
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
	if int(req.BlockIndex) >= len(meta.FileList) {
		return http.StatusBadRequest, nil, errs.NewAPIError(constants.ErrParams, "index out of range")
	}
	output := &codec.StreamInfo{
		Stream: newPartReader(sctx, meta.FileList[req.BlockIndex].FileId),
	}
	return http.StatusOK, output, nil
}
