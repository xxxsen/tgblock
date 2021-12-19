package meta

import (
	"net/http"
	"tgblock/module"
	"tgblock/module/constants"
	"tgblock/processor"

	"github.com/gin-gonic/gin"
)

func GetFileInfo(sctx *module.ServiceContext, ctx *gin.Context, params interface{}) (int, interface{}, error) {
	req := params.(*GetFileInfoRequest)
	if len(req.FileId) == 0 {
		return http.StatusBadRequest, nil, module.NewAPIError(constants.ErrParams, "invalid fileid")
	}
	process := processor.NewFileProcessor(sctx.Bot)
	meta, err := process.GetFileMeta(ctx, req.FileId)
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
