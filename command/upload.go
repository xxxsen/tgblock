package command

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"
	"tgblock/client"
	"tgblock/hasher"
)

func init() {
	Regist(&CmdUpload{})
}

type CmdUpload struct {
	file *string
}

func (c *CmdUpload) Name() string {
	return "upload"
}

func (c *CmdUpload) Args(f *flag.FlagSet) {
	c.file = f.String("file", "", "file path")
}

func (c *CmdUpload) Check() bool {
	if len(*c.file) == 0 {
		log.Printf("file is nil")
		return false
	}
	return true
}

func (c *CmdUpload) Exec(ctx context.Context, cli *client.Client) error {
	stat, err := os.Stat(*c.file)
	if err != nil {
		log.Printf("read file info fail, err:%v", err)
		return nil
	}
	if stat.Size() > cli.MaxFileSize() {
		log.Printf("file size out of limit, max size:%d, file size:%d", cli.MaxFileSize(), stat.Size())
		return nil
	}
	md5str, err := hasher.CalcMD5(*c.file)
	if err != nil {
		log.Printf("calc md5 fail, err:%v", err)
		return err
	}
	reader, err := os.Open(*c.file)
	if err != nil {
		log.Printf("open file for upload fail, err:%v", err)
		return err
	}
	defer reader.Close()
	rsp, err := cli.BlockUpload(ctx, &client.BlockUploadRequest{
		Name:   filepath.Base(*c.file),
		Hash:   md5str,
		Size:   stat.Size(),
		Reader: reader,
	})
	if err != nil {
		log.Printf("upload file fail, err:%v", err)
		return err
	}
	log.Printf("upload succ, fileid:%s", rsp.FileId)
	return nil
}
