package module

import (
	"github.com/gin-gonic/gin"
)

type CommonAuth interface {
	Auth(sctx *ServiceContext, ctx *gin.Context) (bool, error)
}

var NoAuth = &noAuth{}
var TokenAuth = &tokenAuth{}

type noAuth struct {
}

func (a *noAuth) Auth(sctx *ServiceContext, ctx *gin.Context) (bool, error) {
	return true, nil
}

type tokenAuth struct {
}

func (a *tokenAuth) Auth(sctx *ServiceContext, ctx *gin.Context) (bool, error) {
	token := ctx.GetHeader("access_token")
	if token != sctx.AccessToken {
		return false, nil
	}
	return true, nil
}
