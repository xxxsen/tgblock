package codec

import (
	"github.com/gin-gonic/gin"
)

var DefaultNopCodec Codec = &NopCodec{}

type NopCodec struct {
}

func (c *NopCodec) Encode(ctx *gin.Context, code int, input interface{}, err error) error {
	ctx.Status(code)
	return nil
}

func (c *NopCodec) Decode(ctx *gin.Context, output interface{}) error {
	return nil
}
