package upload

import (
	"net/http"
	"tgblock/coder/errs"
	"tgblock/module"
	"tgblock/module/constants"
	"tgblock/processor"

	"github.com/gin-gonic/gin"
)

func BlockUploadBegin(sctx *module.ServiceContext, ctx *gin.Context, params interface{}) (int, interface{}, error) {
	req := params.(*BlockUploadBeginRequest)
	if len(req.Hash) == 0 ||
		len(req.Hash) > 128 || len(req.Name) == 0 ||
		len(req.Name) > 1024 || req.FileSize == 0 ||
		req.FileSize > sctx.MaxFileSize {
		return http.StatusBadRequest, nil, errs.NewAPIError(constants.ErrParams, "invalid params")
	}
	uploader := processor.NewFileProcessor(sctx.Bot)
	begin, err := uploader.CreateFileUpload(ctx, &processor.CreateFileUploadRequest{
		Name:      req.Name,
		FileSize:  req.FileSize,
		BlockSize: sctx.BlockSize,
		HASH:      req.Hash,
	})
	if err != nil {
		return http.StatusInternalServerError, nil, errs.WrapError(constants.ErrUnknown, "call create upload fail", err)
	}
	return http.StatusOK, &BlockUploadBeginResponse{
		UploadId: begin.UploadId,
	}, nil
}
