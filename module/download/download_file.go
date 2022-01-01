package download

import (
	"net/http"

	"github.com/xxxsen/log"
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
	fctxReader := NewFileContextReadSeeker(ctx, sctx, meta)
	//write download info
	output := &codec.StreamInfo{
		Stream: fctxReader,
		Name:   meta.Name,
		Mtime:  meta.CreateTime,
		DeferFunc: func() {
			err := fctxReader.Close()
			log.Debugf("fileid:%s read stream finish, close it, err:%v", fileid, err)
		},
	}
	return http.StatusOK, output, nil
}
