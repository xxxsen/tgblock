package codec

import "net/http"

func MakeCodec(enc Codec, dec Codec) Codec {
	return &combineCodec{
		enc: enc,
		dec: dec,
	}
}

type Encoder interface {
	Encode(request *http.Request, params interface{}) error
}

type Decoder interface {
	Decode(writer *http.Response, target interface{}) error
}

type Codec interface {
	Encoder
	Decoder
}

type combineCodec struct {
	enc Codec
	dec Codec
}

func (c *combineCodec) Encode(request *http.Request, params interface{}) error {
	return c.enc.Encode(request, params)
}

func (c *combineCodec) Decode(rsp *http.Response, target interface{}) error {
	return c.dec.Decode(rsp, target)
}
