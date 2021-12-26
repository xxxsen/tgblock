package meta

import (
	"encoding/base64"
	"net/http"

	"github.com/xxxsen/tgblock/coder/errs"
	"github.com/xxxsen/tgblock/module"
	"github.com/xxxsen/tgblock/module/constants"
	"github.com/xxxsen/tgblock/module/models"
	"github.com/xxxsen/tgblock/shortten"

	"github.com/gin-gonic/gin"
)

func GetFileInfo(sctx *module.ServiceContext, ctx *gin.Context, params interface{}) (int, interface{}, error) {
	req := params.(*models.GetFileInfoRequest)
	if len(req.FileId) == 0 {
		return http.StatusBadRequest, nil, errs.NewAPIError(constants.ErrParams, "invalid fileid")
	}
	fileid, err := shortten.Decode(ctx, req.FileId)
	if err != nil {
		return http.StatusInternalServerError, nil, errs.NewAPIError(constants.ErrUnMarshal, "decode fileid fail")
	}
	process := sctx.Processor
	meta, err := process.CacheGetFileMeta(ctx, fileid)
	if err != nil {
		return http.StatusInternalServerError, nil, errs.WrapError(constants.ErrIO, "get file meta fail", err)
	}
	rsp := &models.GetFileInfoResponse{
		CreateTime: meta.CreateTime,
		FinishTime: meta.FinishTime,
		FileSize:   meta.FileSize,
		Hash:       meta.FileHash,
		BlockSize:  meta.BlockSize,
		BlockCount: meta.BlockCount,
		FileName:   meta.Name,
		BlockHash:  nil,
		FileMode:   meta.FileMode,
	}
	if len(meta.ExtData) != 0 {
		rsp.ExtData = base64.RawStdEncoding.EncodeToString(meta.ExtData)
	}
	for _, item := range meta.FileList {
		rsp.BlockHash = append(rsp.BlockHash, item.Hash)
	}
	return http.StatusOK, rsp, nil
}
