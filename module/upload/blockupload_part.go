package upload

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/xxxsen/tgblock/coder/errs"
	"github.com/xxxsen/tgblock/module"
	"github.com/xxxsen/tgblock/module/constants"
	"github.com/xxxsen/tgblock/processor"

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
		sindex   = ctx.Request.FormValue("index")
	)
	index, err := strconv.ParseInt(sindex, 10, 64)
	if err != nil || index < 0 {
		return http.StatusBadRequest, nil, errs.NewAPIError(constants.ErrParams, "invalid index")
	}
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
	uploader := sctx.Processor
	_, err = uploader.PartFileUpload(ctx, &processor.PartFileUploadRequest{
		UploadId:   uploadid,
		HASH:       hash,
		PartSize:   size,
		Reader:     file,
		BlockIndex: index,
	})
	if err != nil {
		return http.StatusBadRequest, nil, errs.WrapError(constants.ErrIO, "upload fail", err)
	}
	return http.StatusOK, nil, nil
}
