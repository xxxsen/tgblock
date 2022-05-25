package models

type PostUploadResponse struct {
	FileId string `json:"fileid"`
	Hash   string `json:"hash"`
	Size   int64  `json:"size"`
}

type BlockUploadBeginRequest struct {
	FileSize int64 `json:"file_size"`
}

type BlockUploadBeginResponse struct {
	UploadId string `json:"upload_id"`
}

type FileBlock struct {
	FileId    string `json:"file_id"`
	Hash      string `json:"hash"`
	Tagid     uint32 `json:"tag_id"`
	BlockSize int64  `json:"block_size"`
}

type BlockUploadPartResponse struct {
	Block FileBlock `json:"block"`
}

type BlockUploadEndRequest struct {
	UploadId string      `json:"upload_id"`
	Name     string      `json:"name"`
	FileSize int64       `json:"file_size"`
	Hash     string      `json:"hash"`
	FileList []FileBlock `json:"file_list"`
}

type BlockUploadEndResponse struct {
	FileId string `json:"file_id"`
}
