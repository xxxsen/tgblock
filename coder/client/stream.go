package codec

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

var DefaultStreamCodec = NewStreamCodec()

type StreamCodec struct {
}

func NewStreamCodec() *StreamCodec {
	return &StreamCodec{}
}

func (c *StreamCodec) Encode(request *http.Request, params interface{}) error {
	r, ok := params.(io.Reader)
	if !ok {
		return fmt.Errorf("params should be io.Reader")
	}
	request.Body = ioutil.NopCloser(r)
	return nil
}

func (c *StreamCodec) Decode(writer *http.Response, target interface{}) error {
	rc, ok := target.(*io.ReadCloser)
	if !ok {
		return fmt.Errorf("target should be io.ReadCloser PTR")
	}
	*rc = writer.Body
	return nil
}
