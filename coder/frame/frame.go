package frame

type JsonFrame struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func MakeJsonFrame(code int, msg string, data interface{}) *JsonFrame {
	return &JsonFrame{
		Code:    code,
		Message: msg,
		Data:    data,
	}
}

func MakeErrJsonFrame(code int, err string) *JsonFrame {
	return &JsonFrame{
		Code:    code,
		Message: err,
	}
}
