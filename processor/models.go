package processor

import (
	"io"
	"tgblock/protos/gen/tgblock"

	"google.golang.org/protobuf/proto"
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

func FileContextToBytes(fc *tgblock.FileContext) []byte {
	data, _ := proto.Marshal(fc)
	return data
}

func FileContextFromBytes(fc *tgblock.FileContext, data []byte) error {
	if err := proto.Unmarshal(data, fc); err != nil {
		return err
	}
	return nil
}
