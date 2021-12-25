package module

import (
	"fmt"
	"strconv"
	"tgblock/security"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/log"
)

type CommonAuth interface {
	Auth(sctx *ServiceContext, ctx *gin.Context) (bool, error)
}

var NoAuth = &noAuth{}
var SecretAuth = &secretAuth{}

type noAuth struct {
}

func (a *noAuth) Auth(sctx *ServiceContext, ctx *gin.Context) (bool, error) {
	return true, nil
}

type secretAuth struct {
}

func (a *secretAuth) Auth(sctx *ServiceContext, ctx *gin.Context) (bool, error) {
	if len(sctx.SecretId) == 0 {
		log.Errorf("secret id not config, skip auth check")
		return true, nil
	}

	secid := ctx.GetHeader(security.SigSecretId)
	sects := ctx.GetHeader(security.SigSecretTs)
	secsig := ctx.GetHeader(security.SigSecretSig)

	timestamp, err := strconv.ParseInt(sects, 10, 64)
	if err != nil {
		return false, fmt.Errorf("parse timestamp fail, err:%v", err)
	}
	if secid != sctx.SecretId {
		return false, fmt.Errorf("secretid not match")
	}
	if len(sects) == 0 {
		return false, fmt.Errorf("secret_timestamp not found")
	}
	if len(secsig) == 0 {
		return false, fmt.Errorf("secret_sig not found")
	}
	ok, err := security.CheckSig(sctx.SecretId, sctx.SecretKey, secsig, timestamp)
	if err != nil {
		return false, err
	}
	return ok, nil
}
