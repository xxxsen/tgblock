package upload

import (
	"net/http"

	"github.com/xxxsen/tgblock/coder/errs"
	"github.com/xxxsen/tgblock/module"
	"github.com/xxxsen/tgblock/module/constants"
	"github.com/xxxsen/tgblock/module/models"
	"github.com/xxxsen/tgblock/processor"

	"github.com/xxxsen/tgblock/shortten"

	"github.com/gin-gonic/gin"
)

func BlockUploadEnd(sctx *module.ServiceContext, ctx *gin.Context, params interface{}) (int, interface{}, error) {
	req := params.(*models.BlockUploadEndRequest)
	if len(req.UploadId) == 0 {
		return http.StatusBadRequest, nil, errs.NewAPIError(constants.ErrParams, "invalid upload id")
	}
	uploader := sctx.Processor
	finish, err := uploader.FinishFileUpload(ctx, &processor.FinishFileUploadRequest{
		UploadId: req.UploadId,
	})
	if err != nil {
		return http.StatusInternalServerError, nil, errs.WrapError(constants.ErrIO, "call finish upload fail", err)
	}
	fileid, err := shortten.Encode(ctx, finish.FileId)
	if err != nil {
		return http.StatusInternalServerError, nil, errs.NewAPIError(constants.ErrMarshal, "encode fileid fail")
	}
	return http.StatusOK, &models.BlockUploadEndResponse{
		FileId:     fileid,
		CreateTime: finish.CreateTime,
		FinishTime: finish.FinishTime,
		Size:       finish.Size,
		Hash:       finish.Hash,
		BlockSize:  finish.BlockSize,
		BlockCount: finish.BlockCount,
	}, nil
}
