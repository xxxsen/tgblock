package download

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"tgblock/coder/errs"
	codec "tgblock/coder/server"
	"tgblock/module"
	"tgblock/module/constants"
	"tgblock/module/models"
	"tgblock/processor"
	"tgblock/shortten"

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
	proc := processor.NewFileProcessor(sctx.Bot)
	meta, err := proc.GetFileMeta(ctx, fileid)
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
