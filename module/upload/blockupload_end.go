package upload

import (
	"fmt"
	"net/http"

	"github.com/xxxsen/tgblock/coder/errs"
	"github.com/xxxsen/tgblock/module"
	"github.com/xxxsen/tgblock/module/constants"
	"github.com/xxxsen/tgblock/module/models"
	"github.com/xxxsen/tgblock/processor"
	"github.com/xxxsen/tgblock/protos/gen/tgblock"

	"github.com/gin-gonic/gin"
)

func decodeDataFileIds(proc *processor.FileProcessor, filelist []models.FileBlock) error {
	for _, file := range filelist {
		fidctx, err := proc.DecryptFileId(file.FileId)
		if err != nil {
			return err
		}
		if fidctx.FileType != int32(tgblock.FileType_FileType_Data) {
			return errs.NewAPIError(constants.ErrParams, fmt.Sprintf("invalid fileid type, id:%s", file.FileId))
		}
		file.FileId = fidctx.FileId
	}
	return nil
}

func BlockUploadEnd(sctx *module.ServiceContext, ctx *gin.Context, params interface{}) (int, interface{}, error) {
	req := params.(*models.BlockUploadEndRequest)
	if len(req.UploadId) == 0 {
		return http.StatusOK, nil, errs.NewAPIError(constants.ErrParams, "invalid upload id")
	}
	uploader := sctx.Processor

	if err := decodeDataFileIds(uploader, req.FileList); err != nil {
		return http.StatusOK, nil, err
	}

	finish, err := uploader.FinishFileUpload(ctx, &processor.FinishFileUploadRequest{
		UploadId: req.UploadId,
	})
	if err != nil {
		return http.StatusOK, nil, err
	}
	encfileid, err := uploader.EncryptFileId(finish.FileId, int32(tgblock.FileType_FileType_Index))
	if err != nil {
		return http.StatusOK, nil, err
	}

	return http.StatusOK, &models.BlockUploadEndResponse{
		FileId: encfileid,
	}, nil
}
