package codec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"tgblock/coder/frame"
)

var DefaultJsonCodec = NewJsonCodec()

type JsonCodec struct {
}

func NewJsonCodec() *JsonCodec {
	return &JsonCodec{}
}

func (c *JsonCodec) Encode(request *http.Request, params interface{}) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Body = ioutil.NopCloser(bytes.NewReader(data))
	return nil
}

func (c *JsonCodec) Decode(rsp *http.Response, target interface{}) error {
	frame := &frame.JsonFrame{
		Data: target,
	}
	defer rsp.Body.Close()
	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, frame); err != nil {
		return err
	}
	if frame.Code != 0 {
		return fmt.Errorf("logic fail, code:%d, errmsg:%s", frame.Code, frame.Message)
	}
	return nil
}
