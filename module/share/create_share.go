package share

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"tgblock/coder/errs"
	"tgblock/module"
	"tgblock/module/constants"
	"tgblock/module/models"
	"tgblock/protos/gen/tgblock"
	"tgblock/security"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

func CreateShare(sctx *module.ServiceContext, ctx *gin.Context, params interface{}) (int, interface{}, error) {
	req := params.(*models.CreateShareRequest)
	if len(req.FileId) == 0 || len(req.Key) == 0 || len(req.Key) > 32 {
		return http.StatusBadRequest, nil, errs.NewAPIError(constants.ErrParams, "invalid params")
	}
	sblock := &tgblock.ShareBlock{
		Fileid:    []byte(req.FileId),
		Timestamp: req.ExpireTime,
	}
	data, err := proto.Marshal(sblock)
	if err != nil {
		return http.StatusInternalServerError, nil, errs.WrapError(constants.ErrUnknown, "encode fail", err)
	}
	data, err = security.EncryptByKey32(req.Key, data)
	if err != nil {
		return http.StatusInternalServerError, nil, errs.WrapError(constants.ErrUnknown, "encrypt fail", err)
	}
	return http.StatusOK, &models.CreateShareResponse{
		URL: buildURL(sctx, req.Key, base64.RawURLEncoding.EncodeToString(data)),
	}, nil
}

func buildURL(sctx *module.ServiceContext, key string, code string) string {
	return fmt.Sprintf("%s://%s/api/share/getshare?code=%s&key=%s", sctx.Schema, sctx.Domain, code, key)
}
