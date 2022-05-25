package processor

import (
	"io"

	"github.com/xxxsen/tgblock/protos/gen/tgblock"
	"google.golang.org/protobuf/proto"
)

type CreateFileUploadRequest struct {
	FileSize int64
}

type CreateFileUploadResponse struct {
	UploadId string
}

type PartFileUploadRequest struct {
	UploadId string
	Reader   io.Reader
	PartSize int64
}

type PartFileUploadResponse struct {
	FileId    string
	Hash      string
	TagId     uint32
	BlockSize int64
}

type FileBlock struct {
	FileId    string
	Hash      string
	TagId     uint32
	BlockSize int64
}

type FinishFileUploadRequest struct {
	UploadId   string
	FileName   string
	FileIdList []FileBlock
}

type FinishFileUploadResponse struct {
	FileId string
	Hash   string
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
