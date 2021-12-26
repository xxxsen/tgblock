package codec

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/xxxsen/tgblock/coder/errs"
	"github.com/xxxsen/tgblock/coder/frame"

	"github.com/gin-gonic/gin"
)

var DefaultJsonCodec Codec = &JsonCodec{}

type JsonCodec struct {
}

func (c *JsonCodec) Encode(ctx *gin.Context, code int, input interface{}, err error) error {
	if code == http.StatusOK {
		ctx.JSON(http.StatusOK, frame.MakeJsonFrame(0, "", input))
		return nil
	}

	e := errs.AsAPIError(err)
	ctx.JSON(code, frame.MakeErrJsonFrame(e.Code, e.Error()))
	return nil
}

func (c *JsonCodec) Decode(ctx *gin.Context, output interface{}) error {
	data, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, output); err != nil {
		return err
	}
	return nil
}
