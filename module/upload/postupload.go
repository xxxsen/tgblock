package upload

import (
	"fmt"
	"net/http"
	"tgblock/module"
	"tgblock/module/constants"
	"tgblock/processor"
	"tgblock/shortten"

	"github.com/gin-gonic/gin"
)

func PostUpload(sctx *module.ServiceContext, ctx *gin.Context, params interface{}) (int, interface{}, error) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		return http.StatusBadRequest, nil, module.WrapError(constants.ErrParams, "open file invalid", err)
	}
	defer file.Close()
	var (
		name = header.Filename
		size = header.Size
	)
	if size > sctx.BlockSize {
		return http.StatusBadRequest, nil,
			module.NewAPIError(constants.ErrParams, fmt.Sprintf("size exceed, should less than:%d", sctx.BlockSize))
	}
	if len(name) > constants.MaxFileName {
		return http.StatusBadRequest, nil,
			module.NewAPIError(constants.ErrParams, fmt.Sprintf("name too long, should less than:%d", constants.MaxFileName))
	}
	uploader := processor.NewFileProcessor(sctx.Bot)

	begin, err := uploader.CreateFileUpload(ctx, &processor.CreateFileUploadRequest{
		Name:      name,
		FileSize:  size,
		BlockSize: sctx.BlockSize,
	})
	if err != nil {
		return http.StatusInternalServerError, nil, module.WrapError(constants.ErrUnknown, "create upload fail", err)
	}
	part, err := uploader.PartFileUpload(ctx, &processor.PartFileUploadRequest{
		UploadId: begin.UploadId,
		Reader:   file,
		PartSize: size,
	})
	if err != nil {
		return http.StatusInternalServerError, nil, module.WrapError(constants.ErrUnknown, "upload part fail", err)
	}
	finish, err := uploader.FinishFileUpload(ctx, &processor.FinishFileUploadRequest{
		UploadId: begin.UploadId,
	})
	if err != nil {
		return http.StatusInternalServerError, nil, module.WrapError(constants.ErrUnknown, "finish upload failed", err)
	}
	fileid, err := shortten.Encode(ctx, finish.FileId)
	if err != nil {
		return http.StatusInternalServerError, nil, module.NewAPIError(constants.ErrMarshal, "encode fileid fail")
	}
	return http.StatusOK, &PostUploadResponse{
		FileId: fileid,
		Hash:   part.Hash,
		Size:   size,
	}, nil
}
