package codec

import "github.com/gin-gonic/gin"

type Encoder interface {
	Encode(ctx *gin.Context, code int, input interface{}, err error) error
}

type Decoder interface {
	Decode(ctx *gin.Context, output interface{}) error
}

type Codec interface {
	Encoder
	Decoder
}

func MakeCodec(enc Encoder, dec Decoder) Codec {
	return &combineCodec{
		enc: enc,
		dec: dec,
	}
}

func MakeEncoder(enc Encoder) Codec {
	return &combineCodec{
		enc: enc,
		dec: DefaultNopCodec,
	}
}

func MakeDecoder(dec Decoder) Codec {
	return &combineCodec{
		enc: DefaultNopCodec,
		dec: dec,
	}
}

type combineCodec struct {
	enc Encoder
	dec Decoder
}

func (c *combineCodec) Encode(ctx *gin.Context, code int, input interface{}, err error) error {
	return c.enc.Encode(ctx, code, input, err)
}

func (c *combineCodec) Decode(ctx *gin.Context, output interface{}) error {
	return c.dec.Decode(ctx, output)
}
