package download

import (
	"context"
	"io"
	"sync"
	"tgblock/module"
	"tgblock/processor"
)

type multiBlockReader struct {
	meta       *processor.FileContext
	sctx       *module.ServiceContext
	reader     io.Reader
	readerlist []*partReader
}

func newMultiBlockReader(sctx *module.ServiceContext, meta *processor.FileContext) *multiBlockReader {
	mr := &multiBlockReader{
		sctx: sctx,
		meta: meta,
	}
	var readers []io.Reader
	for _, item := range meta.FileList {
		reader := newPartReader(sctx, item.FileId)
		mr.readerlist = append(mr.readerlist, reader)
		readers = append(readers, reader)
	}
	mr.reader = io.MultiReader(readers...)
	return mr
}

func (r *multiBlockReader) Read(buf []byte) (int, error) {
	return r.reader.Read(buf)
}

func (r *multiBlockReader) Close() error {
	var err error
	for _, item := range r.readerlist {
		if e := item.Close(); e != nil {
			err = e
		}
	}
	return err
}

type partReader struct {
	sctx   *module.ServiceContext
	fileid string
	oce    sync.Once
	err    error
	rc     io.ReadCloser
}

func newPartReader(sctx *module.ServiceContext, fileid string) *partReader {
	return &partReader{
		sctx: sctx, fileid: fileid,
	}
}

func (r *partReader) Read(buf []byte) (int, error) {
	r.oce.Do(func() {
		r.rc, r.err = r.sctx.Bot.Download(context.Background(), r.fileid)
	})
	if r.err != nil {
		return 0, r.err
	}
	cnt, err := r.rc.Read(buf)
	if err == io.EOF {
		r.rc.Close()
		r.rc = nil
	}
	return cnt, err
}

func (r *partReader) Close() error {
	if r.rc != nil {
		return r.rc.Close()
	}
	return nil
}
