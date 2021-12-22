package module

import (
	"fmt"
	"strconv"
	"tgblock/security"

	"github.com/gin-gonic/gin"
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
	secid := ctx.GetHeader("secret_id")
	sects := ctx.GetHeader("secret_ts")
	secsig := ctx.GetHeader("secret_sig")

	timestamp, err := strconv.ParseInt(sects, 10, 64)
	if err != nil {
		return false, err
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
