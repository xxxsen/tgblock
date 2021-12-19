package processor

import (
	"encoding/json"
	"io"
)

type CreateFileUploadRequest struct {
	Name      string
	FileSize  int64
	HASH      string
	BlockSize int64
}

type CreateFileUploadResponse struct {
	UploadId string
}

type PartFileUploadRequest struct {
	UploadId string
	Reader   io.Reader
	PartSize int64
	HASH     string
}

type PartFileUploadResponse struct {
	FileId string
	Hash   string
}

type FinishFileUploadRequest struct {
	UploadId string
}

type FinishFileUploadResponse struct {
	FileId     string
	CreateTime int64
	FinishTime int64
	Size       int64
	Hash       string
	BlockSize  int64
	BlockCount int64
}

type FilePart struct {
	FileId string
	Hash   string
}

type FileContext struct {
	Name       string      `json:"name"`
	FileSize   int64       `json:"file_size"`
	FileHash   string      `json:"file_hash"`
	BlockCount int64       `json:"block_count"`
	BlockSize  int64       `json:"block_size"`
	CreateTime int64       `json:"create_time"`
	FinishTime int64       `json:"finish_time"`
	FileList   []*FilePart `json:"file_list"`
}

func (fc *FileContext) ToBytes() []byte {
	data, _ := json.Marshal(fc)
	return data
}

func (fc *FileContext) FromBytes(data []byte) error {
	if err := json.Unmarshal(data, fc); err != nil {
		return err
	}
	return nil
}
