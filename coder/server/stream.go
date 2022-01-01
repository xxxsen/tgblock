package codec

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/xxxsen/tgblock/coder/errs"
	"github.com/xxxsen/tgblock/coder/frame"

	"github.com/gin-gonic/gin"
)

var DefaultStreamCodec Codec = &StreamCodec{}

type StreamCodec struct {
	NopCodec
}

type StreamInfo struct {
	Stream    io.ReadSeeker
	Name      string
	Mtime     int64
	DeferFunc func()
}

func (c *StreamCodec) Encode(ctx *gin.Context, code int, input interface{}, err error) error {
	if code != http.StatusOK {
		e := errs.AsAPIError(err)
		ctx.JSON(code, frame.MakeErrJsonFrame(e.Code, e.Error()))
		return nil
	}
	r, ok := input.(*StreamInfo)
	if !ok {
		return fmt.Errorf("input should be *StreamInfo")
	}
	if r.DeferFunc != nil {
		defer func() {
			r.DeferFunc()
		}()
	}
	if len(r.Name) == 0 {
		r.Name = uuid.NewString()
	}
	http.ServeContent(ctx.Writer, ctx.Request, r.Name, time.Unix(r.Mtime, 0), r.Stream)
	return nil
}
