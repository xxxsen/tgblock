package meta

import (
	"encoding/base64"
	"net/http"
	"tgblock/coder/errs"
	"tgblock/module"
	"tgblock/module/constants"
	"tgblock/module/models"
	"tgblock/processor"
	"tgblock/shortten"

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
	process := processor.NewFileProcessor(sctx.Bot)
	meta, err := process.GetFileMeta(ctx, fileid)
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
