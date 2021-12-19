package meta

import (
	"net/http"
	"tgblock/module"
	"tgblock/module/constants"
	"tgblock/processor"
	"tgblock/shortten"

	"github.com/gin-gonic/gin"
)

func GetFileInfo(sctx *module.ServiceContext, ctx *gin.Context, params interface{}) (int, interface{}, error) {
	req := params.(*GetFileInfoRequest)
	if len(req.FileId) == 0 {
		return http.StatusBadRequest, nil, module.NewAPIError(constants.ErrParams, "invalid fileid")
	}
	fileid, err := shortten.Decode(ctx, req.FileId)
	if err != nil {
		return http.StatusInternalServerError, nil, module.NewAPIError(constants.ErrUnMarshal, "decode fileid fail")
	}
	process := processor.NewFileProcessor(sctx.Bot)
	meta, err := process.GetFileMeta(ctx, fileid)
	if err != nil {
		return http.StatusInternalServerError, nil, module.WrapError(constants.ErrIO, "get file meta fail", err)
	}
	rsp := &GetFileInfoResponse{
		CreateTime: meta.CreateTime,
		FinishTime: meta.FinishTime,
		FileSize:   meta.FileSize,
		Hash:       meta.FileHash,
		BlockSize:  meta.BlockSize,
		BlockCount: meta.BlockCount,
		BlockHash:  nil,
	}
	for _, item := range meta.FileList {
		rsp.BlockHash = append(rsp.BlockHash, item.Hash)
	}
	return http.StatusOK, rsp, nil
}
