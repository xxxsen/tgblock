package upload

import (
	"fmt"
	"net/http"

	"github.com/xxxsen/tgblock/coder/errs"
	"github.com/xxxsen/tgblock/module"
	"github.com/xxxsen/tgblock/module/constants"
	"github.com/xxxsen/tgblock/module/models"
	"github.com/xxxsen/tgblock/processor"

	"github.com/xxxsen/tgblock/shortten"

	"github.com/gin-gonic/gin"
)

func PostUpload(sctx *module.ServiceContext, ctx *gin.Context, params interface{}) (int, interface{}, error) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		return http.StatusBadRequest, nil, errs.WrapError(constants.ErrParams, "open file invalid", err)
	}
	defer file.Close()
	var (
		name = header.Filename
		size = header.Size
	)
	if size > sctx.BlockSize {
		return http.StatusBadRequest, nil,
			errs.NewAPIError(constants.ErrParams, fmt.Sprintf("size exceed, should less than:%d", sctx.BlockSize))
	}
	if len(name) > constants.MaxFileName {
		return http.StatusBadRequest, nil,
			errs.NewAPIError(constants.ErrParams, fmt.Sprintf("name too long, should less than:%d", constants.MaxFileName))
	}
	uploader := sctx.Processor

	begin, err := uploader.CreateFileUpload(ctx, &processor.CreateFileUploadRequest{
		Name:      name,
		FileSize:  size,
		BlockSize: sctx.BlockSize,
	})
	if err != nil {
		return http.StatusInternalServerError, nil, errs.WrapError(constants.ErrUnknown, "create upload fail", err)
	}
	part, err := uploader.PartFileUpload(ctx, &processor.PartFileUploadRequest{
		UploadId: begin.UploadId,
		Reader:   file,
		PartSize: size,
	})
	if err != nil {
		return http.StatusInternalServerError, nil, errs.WrapError(constants.ErrUnknown, "upload part fail", err)
	}
	finish, err := uploader.FinishFileUpload(ctx, &processor.FinishFileUploadRequest{
		UploadId: begin.UploadId,
	})
	if err != nil {
		return http.StatusInternalServerError, nil, errs.WrapError(constants.ErrUnknown, "finish upload failed", err)
	}
	fileid, err := shortten.Encode(ctx, finish.FileId)
	if err != nil {
		return http.StatusInternalServerError, nil, errs.NewAPIError(constants.ErrMarshal, "encode fileid fail")
	}
	return http.StatusOK, &models.PostUploadResponse{
		FileId: fileid,
		Hash:   part.Hash,
		Size:   size,
	}, nil
}
