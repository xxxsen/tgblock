package codec

import (
	"fmt"
	"io"
	"net/http"

	"github.com/xxxsen/tgblock/coder/errs"
	"github.com/xxxsen/tgblock/coder/frame"

	"github.com/gin-gonic/gin"
)

var DefaultStreamCodec Codec = &StreamCodec{}

type StreamCodec struct {
	NopCodec
}

type StreamInfo struct {
	Stream io.ReadCloser
	Size   int64
	Name   string
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
	defer r.Stream.Close()

	if r.Size != 0 {
		ctx.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", r.Size))
	}
	if len(r.Name) != 0 {
		ctx.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", r.Name))
	}
	cnt, err := io.Copy(ctx.Writer, r.Stream)
	if err != nil {
		return err
	}
	if r.Size != 0 && cnt != r.Size {
		return fmt.Errorf("size not match, acquire:%d, got:%d", r.Size, cnt)
	}
	return nil
}
