package upload

import (
	"net/http"
	"tgblock/module"
	"tgblock/module/constants"
	"tgblock/processor"

	"github.com/gin-gonic/gin"
)

func BlockUploadEnd(sctx *module.ServiceContext, ctx *gin.Context, params interface{}) (int, interface{}, error) {
	req := params.(*BlockUploadEndRequest)
	if len(req.UploadId) == 0 {
		return http.StatusBadRequest, nil, module.NewAPIError(constants.ErrParams, "invalid upload id")
	}
	uploader := processor.NewFileProcessor(sctx.Bot)
	finish, err := uploader.FinishFileUpload(ctx, &processor.FinishFileUploadRequest{
		UploadId: req.UploadId,
	})
	if err != nil {
		return http.StatusInternalServerError, nil, module.WrapError(constants.ErrIO, "call finish upload fail", err)
	}
	return http.StatusOK, &BlockUploadEndResponse{
		FileId:     finish.FileId,
		CreateTime: finish.CreateTime,
		FinishTime: finish.FinishTime,
		Size:       finish.Size,
		Hash:       finish.Hash,
		BlockSize:  finish.BlockSize,
		BlockCount: finish.BlockCount,
	}, nil
}
