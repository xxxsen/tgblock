package models

type GetSysInfoRequest struct {
	Timestamp int64 `schema:"timestamp"`
}

type GetSysInfoResponse struct {
	MaxFileSize int64 `json:"max_file_size"`
	BlockSize   int64 `json:"block_size"`
}
