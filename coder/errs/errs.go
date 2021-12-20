package errs

import (
	"fmt"
	"tgblock/module/constants"
)

type APIError struct {
	Code   int    `json:"code"`
	Errmsg string `json:"errmsg"`
	Err    error  `json:"-"`
}

func NewAPIError(code int, errmsg string) *APIError {
	return WrapError(code, errmsg, nil)
}

func WrapError(code int, errmsg string, err error) *APIError {
	return &APIError{
		Code:   code,
		Errmsg: errmsg,
		Err:    err,
	}
}

func (e *APIError) Error() string {
	return fmt.Sprintf("APIERROR:[code:%d, errmsg:%s, err:%v]", e.Code, e.Errmsg, e.Err)
}

func AsAPIError(err error) *APIError {
	e, ok := err.(*APIError)
	if ok {
		return e
	}
	return WrapError(constants.ErrUnknown, "unknow error", err)
}
