package upload

import (
	"fmt"
	"net/http"
	"tgblock/coder/errs"
	"tgblock/module"
	"tgblock/module/constants"
	"tgblock/processor"

	"github.com/gin-gonic/gin"
)

func BlockUploadPart(sctx *module.ServiceContext, ctx *gin.Context, params interface{}) (int, interface{}, error) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		return http.StatusBadRequest, nil, errs.WrapError(constants.ErrParams, "open file invalid", err)
	}
	defer file.Close()
	var (
		size     = header.Size
		uploadid = ctx.Request.FormValue("uploadid")
		hash     = ctx.Request.FormValue("hash")
	)
	if size > sctx.BlockSize {
		return http.StatusBadRequest, nil,
			errs.WrapError(constants.ErrParams, fmt.Sprintf("size exceed, should less than:%d", sctx.BlockSize), err)
	}
	if len(uploadid) == 0 || len(hash) == 0 {
		return http.StatusBadRequest, nil, errs.NewAPIError(constants.ErrParams, "uploadid/hash is nil")
	}
	if size == 0 {
		return http.StatusBadRequest, nil,
			errs.NewAPIError(constants.ErrParams, "size == 0")
	}
	uploader := processor.NewFileProcessor(sctx.Bot)
	_, err = uploader.PartFileUpload(ctx, &processor.PartFileUploadRequest{
		UploadId: uploadid,
		HASH:     hash,
		PartSize: size,
		Reader:   file,
	})
	if err != nil {
		return http.StatusBadRequest, nil, errs.WrapError(constants.ErrIO, "upload fail", err)
	}
	return http.StatusOK, nil, nil
}
