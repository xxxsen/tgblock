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

func DownloadBlock(sctx *module.ServiceContext, ctx *gin.Context, params interface{}) (int, interface{}, error) {
	req := params.(*models.DownloadBlockRequest)
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
	var sz int64
	if meta.GetForceZero() {
		rc = http.NoBody
	} else if len(meta.ExtData) != 0 {
		rc = ioutil.NopCloser(bytes.NewReader(meta.ExtData))
		sz = int64(len(meta.ExtData))
	} else {
		if int(req.BlockIndex) >= len(meta.FileList) {
			return http.StatusBadRequest, nil, errs.NewAPIError(constants.ErrParams, "index out of range")
		}
		rc = newPartReader(sctx, meta.FileList[req.BlockIndex].FileId)
	}
	output := &codec.StreamInfo{
		Stream: rc,
		Size:   sz,
	}
	return http.StatusOK, output, nil
}
