package module

import (
	"encoding/json"
	"io/ioutil"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

type Codec interface {
	Decode(ctx *gin.Context, targetType interface{}) (interface{}, error)
}

type JsonCodec struct {
}

func getTargetInstance(targetType interface{}) interface{} {
	if targetType == nil {
		return nil

	}
	dataType := reflect.TypeOf(targetType).Elem()
	ptr := reflect.New(dataType)
	return ptr.Interface()
}

func (c *JsonCodec) Decode(ctx *gin.Context, targetType interface{}) (interface{}, error) {
	inst := getTargetInstance(targetType)
	if inst == nil {
		return nil, nil
	}
	data, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, inst); err != nil {
		return nil, err
	}
	return inst, nil
}

type URLCodec struct {
}

func (c *URLCodec) Decode(ctx *gin.Context, targetType interface{}) (interface{}, error) {
	inst := getTargetInstance(targetType)
	if inst == nil {
		return nil, nil
	}
	dec := schema.NewDecoder()
	dec.IgnoreUnknownKeys(true)
	if err := dec.Decode(inst, ctx.Request.URL.Query()); err != nil {
		return nil, err
	}
	return inst, nil

}

var (
	DefaultJsonCodec = &JsonCodec{}
	DefaultURLCodec  = &URLCodec{}
)
