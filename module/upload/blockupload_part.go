package upload

import (
	"fmt"
	"net/http"
	"tgblock/module"
	"tgblock/module/constants"
	"tgblock/processor"

	"github.com/gin-gonic/gin"
)

func BlockUploadPart(sctx *module.ServiceContext, ctx *gin.Context, params interface{}) (int, interface{}, error) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		return http.StatusBadRequest, nil, module.WrapError(constants.ErrParams, "open file invalid", err)
	}
	defer file.Close()
	var (
		size     = header.Size
		uploadid = header.Header.Get("uploadid")
		hash     = header.Header.Get("hash")
	)
	if size > sctx.BlockSize {
		return http.StatusBadRequest, nil,
			module.WrapError(constants.ErrParams, fmt.Sprintf("size exceed, should less than:%d", sctx.BlockSize), err)
	}
	if len(uploadid) == 0 || len(hash) == 0 {
		return http.StatusBadRequest, nil, module.NewAPIError(constants.ErrParams, "uploadid/hash is nil")
	}
	uploader := processor.NewFileProcessor(sctx.Bot)
	_, err = uploader.PartFileUpload(ctx, &processor.PartFileUploadRequest{
		UploadId: uploadid,
		HASH:     hash,
		PartSize: size,
		Reader:   file,
	})
	if err != nil {
		return http.StatusBadRequest, nil, module.WrapError(constants.ErrIO, "upload fail", err)
	}
	return http.StatusOK, nil, nil
}
