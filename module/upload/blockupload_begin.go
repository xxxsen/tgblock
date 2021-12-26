package upload

import (
	"net/http"

	"github.com/xxxsen/tgblock/coder/errs"
	"github.com/xxxsen/tgblock/module"
	"github.com/xxxsen/tgblock/module/constants"
	"github.com/xxxsen/tgblock/module/models"
	"github.com/xxxsen/tgblock/processor"

	"github.com/gin-gonic/gin"
)

func BlockUploadBegin(sctx *module.ServiceContext, ctx *gin.Context, params interface{}) (int, interface{}, error) {
	req := params.(*models.BlockUploadBeginRequest)
	if len(req.Hash) == 0 ||
		len(req.Hash) > 128 || len(req.Name) == 0 ||
		len(req.Name) > 1024 || (req.FileSize == 0 && !req.ForceZero) ||
		req.FileSize > sctx.MaxFileSize {
		return http.StatusBadRequest, nil, errs.NewAPIError(constants.ErrParams, "invalid params")
	}
	uploader := sctx.Processor
	begin, err := uploader.CreateFileUpload(ctx, &processor.CreateFileUploadRequest{
		Name:      req.Name,
		FileSize:  req.FileSize,
		BlockSize: sctx.BlockSize,
		HASH:      req.Hash,
		FileMode:  req.FileMode,
		ForceZero: req.ForceZero,
	})
	if err != nil {
		return http.StatusInternalServerError, nil, errs.WrapError(constants.ErrUnknown, "call create upload fail", err)
	}
	return http.StatusOK, &models.BlockUploadBeginResponse{
		UploadId: begin.UploadId,
	}, nil
}
