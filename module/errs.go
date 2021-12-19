package module

import (
	"fmt"
)

type APIError struct {
	Code   int    `json:"code"`
	Errmsg string `json:"errmsg"`
	Err    error  `json:"-"`
}

func NewAPIError(code int, errmsg string) *APIError {
	return &APIError{
		Code:   code,
		Errmsg: errmsg,
	}
}

func WrapError(code int, errmsg string, err error) *APIError {
	return &APIError{
		Code:   code,
		Errmsg: errmsg,
		Err:    err,
	}
}

func (e *APIError) Error() string {
	return fmt.Sprintf("[APIERROR: code:%d, errmsg:%s, err:%v]", e.Code, e.Errmsg, e.Err)
}

func GinResponse(data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"code":   0,
		"errmsg": "success",
		"data":   data,
	}
}

func GinErrResponse(code int, errmsg string, detail error) map[string]interface{} {
	m := make(map[string]interface{})
	m["code"] = code
	m["errmsg"] = errmsg
	if detail != nil {
		m["debug_msg"] = detail.Error()
	}
	return m
}
