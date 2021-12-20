package codec

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

var DefaultURLCodec Codec = &URLCodec{}

type URLCodec struct {
	NopCodec
}

func (c *URLCodec) Decode(ctx *gin.Context, params interface{}) error {
	dec := schema.NewDecoder()
	dec.IgnoreUnknownKeys(true)
	if err := dec.Decode(params, ctx.Request.URL.Query()); err != nil {
		return err
	}
	return nil
}
