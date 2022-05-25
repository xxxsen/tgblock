package constants

const (
	ErrOK        = 0
	ErrUnknown   = 100000
	ErrParams    = 100001
	ErrMarshal   = 100002
	ErrUnMarshal = 100003
	ErrIO        = 100004
	ErrLock      = 100005
)

const (
	MaxFileNameLen = 1024
	BlockSize      = 20 * 1024 * 1024
)
