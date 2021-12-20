package codec

import (
	"fmt"
	"io"
	"net/http"

	"github.com/technoweenie/multipartstreamer"
)

var DefaultFormFileCodec = NewFormFileCodec()

type FormFileCodec struct {
	NopCodec
}

func NewFormFileCodec() *FormFileCodec {
	return &FormFileCodec{}
}

type FormFileInfo struct {
	Name string
	Size int64
	File io.Reader
}

func (c *FormFileCodec) Encode(request *http.Request, params interface{}) error {
	m, ok := params.(map[string]interface{})
	if !ok {
		return fmt.Errorf("params should be map[string]interface{}")
	}
	if len(m) == 0 {
		return fmt.Errorf("no field found")
	}
	fields := make(map[string]string)
	writer := multipartstreamer.New()
	var formkey string
	var formfile *FormFileInfo
	for field, data := range m {
		ffinfo, ok := data.(*FormFileInfo)
		if ok {
			formkey = field
			formfile = ffinfo
			continue
		}
		fields[field] = fmt.Sprintf("%v", data)
	}
	if err := writer.WriteFields(fields); err != nil {
		return err
	}
	if len(formkey) != 0 {
		if err := writer.WriteReader(formkey, formfile.Name, formfile.Size, formfile.File); err != nil {
			return err
		}
	}
	writer.SetupRequest(request)
	return nil
}
