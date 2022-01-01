package processor

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"time"
	"unicode"

	"github.com/xxxsen/tgblock/hasher"
	"github.com/xxxsen/tgblock/protos/gen/tgblock"

	"github.com/google/uuid"
)

const (
	defaultMiniFileLength = 512 //
)

type FileProcessor struct {
	c *Config
}

func NewFileProcessor(opts ...Option) (*FileProcessor, error) {
	c := &Config{}
	for _, opt := range opts {
		opt(c)
	}
	if c.fcache == nil || c.lcker == nil || c.tgbot == nil || len(c.tempdir) == 0 {
		return nil, fmt.Errorf("invalid config")
	}
	dir := c.tempdir + "/fileupload"
	if err := os.RemoveAll(dir); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	c.tempdir = dir
	return &FileProcessor{
		c: c,
	}, nil
}

func (p *FileProcessor) CreateFileUpload(ctx context.Context,
	req *CreateFileUploadRequest) (*CreateFileUploadResponse, error) {

	if len(req.HASH) > 128 || len(req.Name) == 0 ||
		len(req.Name) > 1024 || req.BlockSize == 0 || (req.FileSize == 0 && !req.ForceZero) {

		return nil, fmt.Errorf("invalid params")
	}
	if req.FileMode == 0 {
		req.FileMode = 0755
	}
	fileid := uuid.NewString()
	fctx := &tgblock.FileContext{
		Name:       req.Name,
		FileSize:   req.FileSize,
		FileHash:   req.HASH,
		BlockSize:  req.BlockSize,
		BlockCount: req.FileSize / req.BlockSize,
		CreateTime: time.Now().Unix(),
		FileMode:   req.FileMode,
		ForceZero:  req.ForceZero,
	}
	if req.FileSize%req.BlockSize != 0 {
		fctx.BlockCount += 1
	}
	if err := p.writeToFile(fileid, fctx); err != nil {
		return nil, fmt.Errorf("write ctx to file fail, err:%v", err)
	}
	return &CreateFileUploadResponse{
		UploadId: fileid,
	}, nil
}

func (p *FileProcessor) PartFileUpload(ctx context.Context,
	req *PartFileUploadRequest) (*PartFileUploadResponse, error) {

	var (
		uploadid = req.UploadId
		reader   = hasher.NewMD5Reader(req.Reader)
		hash     = req.HASH
		partsize = req.PartSize
	)

	if !p.isUploadIdValid(uploadid) {
		return nil, fmt.Errorf("invalid upload id")
	}
	lcked := p.c.lcker.Lock(uploadid)
	if !lcked {
		return nil, fmt.Errorf("other process locked")
	}
	defer p.c.lcker.Unlock(uploadid)
	fc, err := p.readFromFile(uploadid)
	if err != nil {
		return nil, fmt.Errorf("read ctx from file fail, err:%v", err)
	}
	if req.BlockIndex < 0 || req.BlockIndex >= fc.BlockCount || req.BlockIndex > int64(len(fc.FileList)) {
		return nil, fmt.Errorf("invalid block count")
	}
	if len(fc.FileList) != int(fc.BlockCount)-1 && req.PartSize > fc.BlockSize {
		return nil, fmt.Errorf("invalid block size, should be:%d", fc.BlockSize)
	}
	//
	var sum string
	var fileid string
	if req.PartSize < defaultMiniFileLength && req.BlockIndex == 0 {
		data, err := ioutil.ReadAll(reader)
		if err != nil {
			return nil, fmt.Errorf("read data fail, err:%v", err)
		}
		sum = reader.GetSum()
		fc.ExtData = data
	} else {
		fileid, err = p.c.tgbot.Upload(ctx, partsize, reader)
		if err != nil {
			return nil, fmt.Errorf("upload fail, err:%v", err)
		}
		sum = reader.GetSum()
	}
	blockInfo := &tgblock.FilePart{
		FileId: fileid,
		Hash:   sum,
	}
	if int(req.BlockIndex) == len(fc.FileList) {
		fc.FileList = append(fc.FileList, blockInfo)
	} else {
		fc.FileList[req.BlockIndex] = blockInfo
	}
	if len(hash) != 0 && sum != hash {
		return nil, fmt.Errorf("checksum not match, should be:%s, get:%s", sum, hash)
	}

	if err := p.writeToFile(uploadid, fc); err != nil {
		return nil, fmt.Errorf("write ctx to file fail")
	}
	return &PartFileUploadResponse{
		Hash: sum,
	}, nil
}

func (p *FileProcessor) FinishFileUpload(ctx context.Context,
	req *FinishFileUploadRequest) (*FinishFileUploadResponse, error) {

	var (
		uploadid = req.UploadId
	)
	if !p.isUploadIdValid(uploadid) {
		return nil, fmt.Errorf("invalid upload id")
	}
	lcked := p.c.lcker.Lock(uploadid)
	if !lcked {
		return nil, fmt.Errorf("other process locked")
	}
	defer p.c.lcker.Unlock(uploadid)
	fc, err := p.readFromFile(uploadid)
	if err != nil {
		return nil, fmt.Errorf("read ctx from file fail, err:%v", err)
	}
	if !fc.ForceZero && fc.FileSize == 0 {
		return nil, fmt.Errorf("zero size file found")
	}
	fc.FinishTime = time.Now().Unix()
	if len(fc.FileHash) == 0 && len(fc.FileList) == 1 {
		fc.FileHash = fc.FileList[0].Hash
	}

	if err := p.checkFileContextValid(fc); err != nil {
		return nil, fmt.Errorf("check file ctx fail, err:%v", err)
	}
	data := FileContextToBytes(fc)
	fileid, err := p.c.tgbot.Upload(ctx, int64(len(data)), bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("finish upload fail, err:%v", err)
	}
	_ = p.removeFile(uploadid)
	return &FinishFileUploadResponse{
		FileId:     fileid,
		CreateTime: fc.CreateTime,
		FinishTime: fc.FinishTime,
		Size:       fc.FileSize,
		Hash:       fc.FileHash,
		BlockSize:  fc.BlockSize,
		BlockCount: fc.BlockCount,
	}, nil
}

func (p *FileProcessor) checkFileContextValid(fc *tgblock.FileContext) error {
	total := len(fc.FileList)
	left := (total - 1) * int(fc.BlockSize)
	right := total * int(fc.BlockSize)
	if int(fc.FileSize) < left || int(fc.FileSize) > right {
		return fmt.Errorf("invalid block num")
	}
	if len(fc.FileList) != int(fc.BlockCount) {
		return fmt.Errorf("invalid block count")
	}
	return nil
}

func (p *FileProcessor) isUploadIdValid(file string) bool {
	if len(file) == 0 {
		return false
	}
	for _, c := range file {
		if !(unicode.IsLetter(c) || unicode.IsDigit(c) || c == rune('-')) {
			return false
		}
	}
	return true
}

func (p *FileProcessor) readFromFile(uploadid string) (*tgblock.FileContext, error) {
	file := p.buildSavePath(uploadid)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	fc := &tgblock.FileContext{}
	if err := FileContextFromBytes(fc, data); err != nil {
		return nil, err
	}
	return fc, nil
}

func (p *FileProcessor) writeToFile(uploadid string, fc *tgblock.FileContext) error {
	file := p.buildSavePath(uploadid)

	if err := ioutil.WriteFile(file, FileContextToBytes(fc), 0644); err != nil {
		return err
	}
	return nil
}

func (p *FileProcessor) removeFile(uploadid string) error {
	file := p.buildSavePath(uploadid)
	return os.Remove(file)
}

func (p *FileProcessor) buildSavePath(uploadid string) string {
	return p.c.tempdir + uploadid
}

func (p *FileProcessor) CacheGetFileMeta(ctx context.Context, fileid string) (*tgblock.FileContext, error) {
	v, err := p.c.fcache.Get(ctx, fileid)
	if err != nil {
		return nil, err
	}
	return v.(*tgblock.FileContext), nil
}

func (p *FileProcessor) GetFileMeta(ctx context.Context, fileid string) (*tgblock.FileContext, error) {
	data, err := p.GetFileData(ctx, fileid)
	if err != nil {
		return nil, err
	}
	fc := &tgblock.FileContext{}
	if err := FileContextFromBytes(fc, data); err != nil {
		return nil, err
	}
	return fc, nil
}

func (p *FileProcessor) GetFileData(ctx context.Context, fileid string) ([]byte, error) {
	reader, err := p.DownloadFile(ctx, fileid)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (p *FileProcessor) DownloadFile(ctx context.Context, fileid string) (io.ReadCloser, error) {
	return p.c.tgbot.DownloadAt(ctx, fileid, 0)
}

func (p *FileProcessor) DownloadFileAt(ctx context.Context, fileid string, index int64) (io.ReadCloser, error) {
	return p.c.tgbot.DownloadAt(ctx, fileid, index)
}
