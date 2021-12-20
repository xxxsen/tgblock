package codec

import (
	"net/http"
	"net/url"

	"github.com/gorilla/schema"
)

var DefaultURLCodec Codec = NewURLCodec()

type URLCodec struct {
	NopCodec
}

func NewURLCodec() *URLCodec {
	return &URLCodec{}
}

func (c *URLCodec) Encode(request *http.Request, params interface{}) error {
	query := url.Values{}
	enc := schema.NewEncoder()
	if err := enc.Encode(params, query); err != nil {
		return err
	}
	request.URL.RawQuery = query.Encode()
	return nil
}
