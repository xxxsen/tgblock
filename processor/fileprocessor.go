package processor

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"math/rand"

	"time"

	"github.com/xxxsen/tgblock/coder/errs"
	"github.com/xxxsen/tgblock/hasher"
	"github.com/xxxsen/tgblock/module/constants"
	"github.com/xxxsen/tgblock/protos/gen/tgblock"
	"github.com/xxxsen/tgblock/security"
	"google.golang.org/protobuf/proto"
)

type FileProcessor struct {
	c *Config
}

func NewFileProcessor(opts ...Option) (*FileProcessor, error) {
	c := &Config{}
	for _, opt := range opts {
		opt(c)
	}
	if c.fcache == nil || c.lcker == nil || c.tgbot == nil || len(c.seckey) == 0 {
		return nil, fmt.Errorf("invalid config")
	}
	return &FileProcessor{
		c: c,
	}, nil
}

func (p *FileProcessor) CreateFileUpload(ctx context.Context,
	req *CreateFileUploadRequest) (*CreateFileUploadResponse, error) {
	fileid, err := p.CreateUploadId(req.FileSize)
	if err != nil {
		return nil, fmt.Errorf("create upload id fail, err:%v", err)
	}
	return &CreateFileUploadResponse{
		UploadId: fileid,
	}, nil
}

func (p *FileProcessor) CreateUploadId(filesize int64) (string, error) {
	idblk := &tgblock.UploadIdContext{
		RandId:    rand.Int63(),
		Timestamp: time.Now().Unix(),
		FileSize:  filesize,
	}
	raw, err := proto.Marshal(idblk)
	if err != nil {
		return "", fmt.Errorf("encode id fail, err:%v", err)
	}
	out, err := security.EncryptByKey32(p.c.seckey, raw)
	if err != nil {
		return "", fmt.Errorf("encrypt id fail, err:%v", err)
	}
	return hex.EncodeToString(out), nil
}

func (p *FileProcessor) DecodeUploadId(id string) (*tgblock.UploadIdContext, error) {
	raw, err := hex.DecodeString(id)
	if err != nil {
		return nil, fmt.Errorf("decode id fail, err:%v", err)
	}
	out, err := security.DecryptByKey32(p.c.seckey, raw)
	if err != nil {
		return nil, fmt.Errorf("decrypt id fail, err:%v", err)
	}
	blk := &tgblock.UploadIdContext{}
	if err := proto.Unmarshal(out, blk); err != nil {
		return nil, fmt.Errorf("unmarshal blk fail, err:%v", err)
	}
	return blk, nil
}

func (p *FileProcessor) BuildTagId(idblk *tgblock.UploadIdContext, fileid string, hash string, size int64) uint32 {
	crc := crc32.NewIEEE()
	crc.Write([]byte(fmt.Sprintf("%d", idblk.RandId)))
	crc.Write([]byte(hash))
	crc.Write([]byte(fileid))
	crc.Write([]byte(fmt.Sprintf("%d", size)))
	return crc.Sum32()
}

func (p *FileProcessor) PartFileUpload(ctx context.Context,
	req *PartFileUploadRequest) (*PartFileUploadResponse, error) {

	var (
		uploadid = req.UploadId
		reader   = hasher.NewMD5Reader(req.Reader)
		partsize = req.PartSize
	)

	if partsize > constants.BlockSize {
		return nil, fmt.Errorf("invalid block size:%d", partsize)
	}

	blk, err := p.DecodeUploadId(uploadid)
	if err != nil {
		return nil, fmt.Errorf("decode upload id fail, err:%v", err)
	}

	lcked := p.c.lcker.Lock(uploadid)
	if !lcked {
		return nil, fmt.Errorf("other process locked")
	}
	defer p.c.lcker.Unlock(uploadid)

	fileid, err := p.c.tgbot.Upload(ctx, partsize, reader)
	if err != nil {
		return nil, fmt.Errorf("upload fail, err:%v", err)
	}
	if reader.GetSize() != int(partsize) {
		return nil, fmt.Errorf("part size not match, get:%d", reader.GetSize())
	}
	sum := reader.GetSum()
	tagid := p.BuildTagId(blk, fileid, sum, partsize)
	return &PartFileUploadResponse{
		Hash:      sum,
		TagId:     tagid,
		FileId:    fileid,
		BlockSize: partsize,
	}, nil
}

func (p *FileProcessor) checkFileBlockListValid(blk *tgblock.UploadIdContext, blocklist []FileBlock) error {
	var fullSize int64
	for idx, block := range blocklist {
		tagid := p.BuildTagId(blk, block.FileId, block.Hash, block.BlockSize)
		if tagid != block.TagId {
			return fmt.Errorf("idx:%d, tagid:%d, blocktagid:%d not match", idx, tagid, block.TagId)
		}
		if block.BlockSize < constants.BlockSize && idx != len(blocklist)-1 {
			return fmt.Errorf("part of block use invalid file size:%d", block.BlockSize)
		}
		fullSize += block.BlockSize
	}
	if fullSize != blk.FileSize {
		return fmt.Errorf("file size not match, get:%d, request:%d", fullSize, blk.FileSize)
	}
	return nil
}

func (p *FileProcessor) buildBlockListHash(blocklist []FileBlock) string {
	hasher := md5.New()
	for _, block := range blocklist {
		hasher.Write([]byte(block.Hash))
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

func (p *FileProcessor) buildFileIdList(blocklist []FileBlock) []string {
	rs := make([]string, 0, len(blocklist))
	for _, block := range blocklist {
		rs = append(rs, block.FileId)
	}
	return rs
}

func (p *FileProcessor) FinishFileUpload(ctx context.Context,
	req *FinishFileUploadRequest) (*FinishFileUploadResponse, error) {

	var (
		uploadid = req.UploadId
		name     = req.FileName
	)
	if len(name) > constants.MaxFileNameLen {
		return nil, errs.NewAPIError(constants.ErrParams, fmt.Sprintf("file name too long, size:%d", len(name)))
	}

	blk, err := p.DecodeUploadId(uploadid)
	if err != nil {
		return nil, errs.WrapError(constants.ErrParams, "decode upload id fail", err)
	}
	if err := p.checkFileBlockListValid(blk, req.FileIdList); err != nil {
		return nil, errs.WrapError(constants.ErrParams, "check block fail", err)
	}

	lcked := p.c.lcker.Lock(uploadid)
	if !lcked {
		return nil, errs.NewAPIError(constants.ErrLock, "other process locked")
	}
	defer p.c.lcker.Unlock(uploadid)

	fc := &tgblock.FileContext{
		Name:       req.FileName,
		FileSize:   blk.FileSize,
		FileHash:   p.buildBlockListHash(req.FileIdList),
		CreateTime: 0,
		FileIds:    p.buildFileIdList(req.FileIdList),
	}

	data := FileContextToBytes(fc)
	fileid, err := p.c.tgbot.Upload(ctx, int64(len(data)), bytes.NewReader(data))
	if err != nil {
		return nil, errs.WrapError(constants.ErrIO, "finish upload fail", err)
	}
	return &FinishFileUploadResponse{
		FileId: fileid,
		Hash:   fc.FileHash,
	}, nil
}

func (p *FileProcessor) EncryptFileId(fileid string, ftype int32) (string, error) {
	fidctx := &tgblock.FileIdContext{
		FileType: ftype,
		FileId:   fileid,
	}
	raw, err := proto.Marshal(fidctx)
	if err != nil {
		return "", errs.WrapError(constants.ErrUnknown, "encode pb fail, err:%v", err)
	}
	out, err := security.EncryptByKey32(p.c.seckey, raw)
	if err != nil {
		return "", errs.WrapError(constants.ErrUnknown, "encrypt fail", err)
	}
	return base64.URLEncoding.EncodeToString(out), nil
}

func (p *FileProcessor) DecryptFileId(fileid string) (*tgblock.FileIdContext, error) {
	raw, err := base64.URLEncoding.DecodeString(fileid)
	if err != nil {
		return nil, errs.WrapError(constants.ErrParams, "decode wrap fail", err)
	}
	out, err := security.DecryptByKey32(p.c.seckey, raw)
	if err != nil {
		return nil, errs.WrapError(constants.ErrParams, "decrypt fail", err)
	}
	fidctx := &tgblock.FileIdContext{}
	if err := proto.Unmarshal(out, fidctx); err != nil {
		return nil, errs.WrapError(constants.ErrParams, "final decode fail", err)
	}
	return nil, err
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
