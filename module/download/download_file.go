package download

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	codec "github.com/xxxsen/tgblock/coder/server"

	"github.com/xxxsen/tgblock/coder/errs"
	"github.com/xxxsen/tgblock/module"
	"github.com/xxxsen/tgblock/module/constants"
	"github.com/xxxsen/tgblock/module/models"
	"github.com/xxxsen/tgblock/shortten"

	"github.com/gin-gonic/gin"
)

func DownloadFile(sctx *module.ServiceContext, ctx *gin.Context, params interface{}) (int, interface{}, error) {
	req := params.(*models.DownloadFileRequest)
	if len(req.FileId) == 0 {
		return http.StatusBadRequest, nil, errs.NewAPIError(constants.ErrParams, "invalid fileid")
	}
	fileid, err := shortten.Decode(ctx, req.FileId)
	if err != nil {
		return http.StatusInternalServerError, nil, errs.WrapError(constants.ErrUnMarshal, "decode fileid fail", err)
	}
	proc := sctx.Processor
	meta, err := proc.CacheGetFileMeta(ctx, fileid)
	if err != nil {
		return http.StatusInternalServerError, nil, errs.WrapError(constants.ErrIO, "read file meta fail", err)
	}
	var rc io.ReadCloser
	if meta.ForceZero {
		rc = http.NoBody
	} else if len(meta.ExtData) > 0 {
		rc = ioutil.NopCloser(bytes.NewReader(meta.ExtData))
	} else {
		rc = newMultiBlockReader(sctx, meta)
	}

	//write download info
	output := &codec.StreamInfo{
		Stream: rc,
		Size:   meta.FileSize,
		Name:   meta.Name,
	}
	return http.StatusOK, output, nil
}
