package models

type PostUploadResponse struct {
	FileId string `json:"fileid"`
	Hash   string `json:"hash"`
	Size   int64  `json:"size"`
}

type BlockUploadBeginRequest struct {
	Name      string `json:"name"`
	FileSize  int64  `json:"file_size"`
	Hash      string `json:"hash"`
	FileMode  int64  `json:"file_mode"`
	ForceZero bool   `json:"force_zero"`
}

type BlockUploadBeginResponse struct {
	UploadId string `json:"upload_id"`
}

type BlockUploadEndRequest struct {
	UploadId string `json:"upload_id"`
}

type BlockUploadEndResponse struct {
	FileId     string `json:"file_id"`
	CreateTime int64  `json:"create_time"`
	FinishTime int64  `json:"finish_time"`
	Size       int64  `json:"size"`
	Hash       string `json:"hash"`
	BlockSize  int64  `json:"block_size"`
	BlockCount int64  `json:"block_count"`
}
