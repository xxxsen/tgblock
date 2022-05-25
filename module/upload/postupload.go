package upload

import (
	"fmt"
	"net/http"

	"github.com/xxxsen/tgblock/coder/errs"
	"github.com/xxxsen/tgblock/module"
	"github.com/xxxsen/tgblock/module/constants"
	"github.com/xxxsen/tgblock/module/models"
	"github.com/xxxsen/tgblock/processor"
	"github.com/xxxsen/tgblock/protos/gen/tgblock"

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
	if size > constants.BlockSize {
		return http.StatusBadRequest, nil,
			errs.NewAPIError(constants.ErrParams, fmt.Sprintf("size exceed, should less than:%d", constants.BlockSize))
	}
	uploader := sctx.Processor
	begin, err := uploader.CreateFileUpload(ctx, &processor.CreateFileUploadRequest{
		FileSize: size,
	})
	if err != nil {
		return http.StatusOK, nil, err
	}
	part, err := uploader.PartFileUpload(ctx, &processor.PartFileUploadRequest{
		UploadId: begin.UploadId,
		Reader:   file,
		PartSize: size,
	})
	if err != nil {
		return http.StatusOK, nil, err
	}
	finish, err := uploader.FinishFileUpload(ctx, &processor.FinishFileUploadRequest{
		UploadId: begin.UploadId,
		FileName: name,
		FileIdList: []processor.FileBlock{
			{
				FileId:    part.FileId,
				Hash:      part.Hash,
				TagId:     part.TagId,
				BlockSize: part.BlockSize,
			},
		},
	})
	if err != nil {
		return http.StatusOK, nil, err
	}
	encfileid, err := uploader.EncryptFileId(finish.FileId, int32(tgblock.FileType_FileType_Index))
	if err != nil {
		return http.StatusOK, nil, err
	}
	return http.StatusOK, &models.PostUploadResponse{
		FileId: encfileid,
		Hash:   finish.Hash,
		Size:   size,
	}, nil
}
