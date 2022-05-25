package upload

import (
	"fmt"
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

	if req.FileSize == 0 || req.FileSize > sctx.MaxFileSize {
		return http.StatusBadRequest, nil, errs.NewAPIError(constants.ErrParams, fmt.Sprintf("size invalid, max:%d", sctx.MaxFileSize))
	}
	proc := sctx.Processor
	rsp, err := proc.CreateFileUpload(ctx, &processor.CreateFileUploadRequest{
		FileSize: req.FileSize,
	})
	if err != nil {
		return http.StatusBadGateway, nil, errs.WrapError(constants.ErrUnknown, "create upload fail", err)
	}
	return http.StatusOK, &models.BlockUploadBeginResponse{
		UploadId: rsp.UploadId,
	}, nil
}
