package upload

import (
	"net/http"
	"tgblock/coder/errs"
	"tgblock/module"
	"tgblock/module/constants"
	"tgblock/module/models"
	"tgblock/processor"
	"tgblock/shortten"

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
