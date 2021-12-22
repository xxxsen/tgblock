package share

import (
	"encoding/base64"
	"net/http"
	"tgblock/coder/errs"
	"tgblock/module"
	"tgblock/module/constants"
	"tgblock/module/download"
	"tgblock/module/models"
	"tgblock/protos/gen/tgblock"
	"tgblock/security"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

func GetShare(sctx *module.ServiceContext, ctx *gin.Context, params interface{}) (int, interface{}, error) {
	req := params.(*models.GetShareRequest)
	if len(req.Code) == 0 || len(req.Key) == 0 {
		return http.StatusBadRequest, nil, errs.NewAPIError(constants.ErrParams, "invalid params")
	}

	raw, err := base64.RawURLEncoding.DecodeString(req.Code)
	if err != nil {
		return http.StatusBadRequest, nil, errs.WrapError(constants.ErrParams, "invalid params", err)
	}
	data, err := security.DecryptByKey32(req.Key, raw)
	if err != nil {
		return http.StatusBadRequest, nil, errs.WrapError(constants.ErrParams, "decrypt fail", err)
	}

	sblock := &tgblock.ShareBlock{}
	if err := proto.Unmarshal(data, sblock); err != nil {
		return http.StatusBadRequest, nil, errs.WrapError(constants.ErrParams, "decode fail", err)
	}
	if sblock.Timestamp != 0 && sblock.Timestamp < time.Now().Unix() {
		return http.StatusBadRequest, nil, errs.NewAPIError(constants.ErrParams, "expired")
	}
	fileid := string(sblock.Fileid)
	return download.DownloadFile(sctx, ctx, &models.DownloadFileRequest{
		FileId: fileid,
	})
}
