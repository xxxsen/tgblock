package models

type GetFileInfoRequest struct {
	FileId string `schema:"file_id"`
}

type GetFileInfoResponse struct {
	CreateTime int64    `json:"create_time"`
	FinishTime int64    `json:"finish_time"`
	FileSize   int64    `json:"file_size"`
	Hash       string   `json:"hash"`
	BlockSize  int64    `json:"block_size"`
	BlockCount int64    `json:"block_count"`
	BlockHash  []string `json:"block_hash"`
	FileName   string   `json:"file_name"`
	FileMode   int64    `json:"file_mode"`
	ExtData    string   `json:"ext_data"`
}
