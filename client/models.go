package client

import "io"

type FileInfo struct {
	Name string
	Size int64
	Hash string
	File io.Reader
}

type BlockUploadRequest struct {
	Name      string
	Size      int64
	Reader    io.Reader
	Hash      string
	Mode      int64
	ForceZero bool
}

type BlockUploadResponse struct {
	FileId string
}
