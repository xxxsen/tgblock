package share

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/xxxsen/tgblock/coder/errs"
	"github.com/xxxsen/tgblock/module"
	"github.com/xxxsen/tgblock/module/constants"
	"github.com/xxxsen/tgblock/module/download"
	"github.com/xxxsen/tgblock/module/models"
	"github.com/xxxsen/tgblock/protos/gen/tgblock"
	"github.com/xxxsen/tgblock/security"

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
