package upload

import (
	"net/http"

	"github.com/xxxsen/tgblock/coder/errs"
	"github.com/xxxsen/tgblock/module"
	"github.com/xxxsen/tgblock/module/constants"
	"github.com/xxxsen/tgblock/module/models"
	"github.com/xxxsen/tgblock/processor"
	"github.com/xxxsen/tgblock/protos/gen/tgblock"

	"github.com/gin-gonic/gin"
)

func BlockUploadPart(sctx *module.ServiceContext, ctx *gin.Context, params interface{}) (int, interface{}, error) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		return http.StatusOK, nil, errs.WrapError(constants.ErrParams, "open file invalid", err)
	}
	defer file.Close()
	var (
		size     = header.Size
		uploadid = ctx.Request.FormValue("uploadid")
	)
	uploader := sctx.Processor
	part, err := uploader.PartFileUpload(ctx, &processor.PartFileUploadRequest{
		UploadId: uploadid,
		PartSize: size,
		Reader:   file,
	})
	if err != nil {
		return http.StatusOK, nil, err
	}
	encfileid, err := uploader.EncryptFileId(part.FileId, int32(tgblock.FileType_FileType_Data))
	if err != nil {
		return http.StatusOK, nil, err
	}
	rsp := &models.BlockUploadPartResponse{
		Block: models.FileBlock{
			FileId:    encfileid,
			Hash:      part.Hash,
			Tagid:     part.TagId,
			BlockSize: part.BlockSize,
		},
	}
	return http.StatusOK, rsp, nil
}
