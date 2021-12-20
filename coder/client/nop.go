package codec

import "net/http"

var DefaultNopCodec = NewNopCodec()

type NopCodec struct {
}

func NewNopCodec() *NopCodec {
	return &NopCodec{}
}

func (c *NopCodec) Encode(request *http.Request, params interface{}) error {
	return nil
}

func (c *NopCodec) Decode(writer *http.Response, target interface{}) error {
	writer.Body.Close()
	return nil
}
