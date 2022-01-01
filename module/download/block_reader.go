package download

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/log"
	"github.com/xxxsen/tgblock/module"
	"github.com/xxxsen/tgblock/utils"

	"github.com/xxxsen/tgblock/protos/gen/tgblock"
)

type ReadSeekCloser interface {
	io.Reader
	io.Closer
	io.Seeker
}

type simpleFileCtx struct {
	*bytes.Reader
}

func (c *simpleFileCtx) Close() error {
	return nil
}

type FileContextReadSeeker struct {
	ctx       *gin.Context
	meta      *tgblock.FileContext
	sctx      *module.ServiceContext
	reader    io.ReadCloser
	blockId   int64
	readIndex int64
}

func NewFileContextReadSeeker(ctx *gin.Context, sctx *module.ServiceContext, meta *tgblock.FileContext) ReadSeekCloser {
	if meta.ForceZero || len(meta.ExtData) > 0 {
		return &simpleFileCtx{Reader: bytes.NewReader(meta.ExtData)}
	}

	return &FileContextReadSeeker{
		ctx:  ctx,
		sctx: sctx,
		meta: meta,
	}
}

func (f *FileContextReadSeeker) switchNextBlock() error {
	err := f.resetReader()
	if err != nil {
		return err
	}
	f.blockId++
	if f.blockId == f.meta.BlockCount {
		return io.EOF
	}
	reader, err := f.blockIdIndexToReader(f.blockId, 0)
	if err != nil {
		return err
	}
	f.reader = reader
	return nil
}

func (f *FileContextReadSeeker) Read(buf []byte) (int, error) {
	cnt, err := f.reader.Read(buf)
	if err == io.EOF {
		err = f.switchNextBlock()
	}
	if cnt > 0 {
		f.readIndex += int64(cnt)
	}
	return cnt, err
}

func (f *FileContextReadSeeker) resetReader() error {
	if f.reader != nil {
		return f.reader.Close()
	}
	return nil
}

func (f *FileContextReadSeeker) Close() error {
	f.blockId = 0
	f.readIndex = 0
	return f.resetReader()
}

func (f *FileContextReadSeeker) blockIdIndexToReader(blockid int64, readindex int64) (*FileContextBlockReadSeeker, error) {
	fileid := f.meta.FileList[blockid].GetFileId()
	fullsize, err := utils.CalcBlockSizeByIndex(f.meta, blockid)
	if err != nil {
		return nil, err
	}
	reader := NewFileContextBlockReader(f.ctx, f.sctx, fileid, fullsize)
	if readindex != 0 {
		_, err = reader.Seek(readindex, io.SeekStart)
		if err != nil {
			reader.Close()
			return nil, err
		}
	}
	return reader, nil
}

func (f *FileContextReadSeeker) Seek(offset int64, whence int) (int64, error) {
	if err := f.resetReader(); err != nil {
		return 0, err
	}

	//TODO:next start 计算错误
	nextStart, err := utils.CalcSeek(f.meta.FileSize, f.readIndex, offset, whence)
	if err != nil {
		return 0, err
	}
	blockid, readindex, err := utils.ReadIndexToBlockIndexOffset(f.meta, nextStart)
	if err != nil {
		return 0, err
	}
	reader, err := f.blockIdIndexToReader(blockid, readindex)
	if err != nil {
		return 0, err
	}
	f.reader = reader
	f.blockId = blockid
	f.readIndex = nextStart
	log.Debugf("file seek, hash:%s, blockid:%d, final readindex:%d, nextstart:%d", f.meta.FileHash, f.blockId, f.readIndex, nextStart)
	return f.readIndex, nil
}

type FileContextBlockReadSeeker struct {
	ctx       *gin.Context
	sctx      *module.ServiceContext
	fileid    string
	reader    io.ReadCloser
	fullSize  int64
	readIndex int64
}

func NewFileContextBlockReader(ctx *gin.Context, sctx *module.ServiceContext, fileid string, fullsize int64) *FileContextBlockReadSeeker {
	return &FileContextBlockReadSeeker{
		ctx:      ctx,
		sctx:     sctx,
		fileid:   fileid,
		fullSize: fullsize,
	}
}

func (f *FileContextBlockReadSeeker) tryOpenStream() error {
	if f.reader != nil {
		return nil
	}
	rc, err := f.sctx.Processor.DownloadFileAt(context.Background(), f.fileid, f.readIndex)
	log.Debugf("fileid:%s, open stream at loc:%d, err:%v", f.fileid, f.readIndex, err)
	if err != nil {
		return err
	}
	f.reader = rc
	return nil
}

func (f *FileContextBlockReadSeeker) Read(buf []byte) (int, error) {
	if err := f.tryOpenStream(); err != nil {
		return 0, fmt.Errorf("open stream at index:%d fail, err:%v", f.readIndex, err)
	}
	cnt, err := f.reader.Read(buf)
	if cnt > 0 {
		f.readIndex += int64(cnt)
	}
	if err == io.EOF {
		e := f.resetReader()
		if e != nil {
			err = e
		}
	}
	return cnt, err
}

func (f *FileContextBlockReadSeeker) resetReader() error {
	var err error
	if f.reader != nil {
		err = f.reader.Close()
		f.reader = nil
	}
	return err
}

func (f *FileContextBlockReadSeeker) Close() error {
	f.readIndex = 0
	return f.resetReader()
}

func (f *FileContextBlockReadSeeker) Seek(offset int64, whence int) (int64, error) {
	if err := f.resetReader(); err != nil {
		return 0, err
	}
	nextStart, err := utils.CalcSeek(f.fullSize, f.readIndex, offset, whence)
	log.Debugf("fileid:%s, reset index from:%d to:%d", f.fileid, f.readIndex, nextStart)
	if err != nil {
		return 0, err
	}
	f.readIndex = nextStart
	return 0, nil
}
