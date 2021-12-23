package command

import (
	"bytes"
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"tgblock/client"
	"tgblock/module/models"
)

func init() {
	Regist(&CmdCat{})
}

type CmdCat struct {
	fileid       *string
	maxSizeLimit *int64
}

func (c *CmdCat) Name() string {
	return "cat"
}

func (c *CmdCat) Args(f *flag.FlagSet) {
	c.fileid = f.String("fileid", "", "file id")
	c.maxSizeLimit = f.Int64("max_size_limit", 4*1024, "max size to output")
}

func (c *CmdCat) Check() bool {
	if len(*c.fileid) == 0 {
		log.Printf("fileid is nil")
		return false
	}
	return true
}

func (c *CmdCat) Exec(ctx context.Context, cli *client.Client) error {
	info, err := cli.GetFileInfo(ctx, &models.GetFileInfoRequest{
		FileId: *c.fileid,
	})
	if err != nil {
		return err
	}
	if info.FileSize > *c.maxSizeLimit {
		log.Printf("file too big to print, size:%d, limit:%d", info.FileSize, *c.maxSizeLimit)
		return nil
	}
	var rc io.ReadCloser
	if len(info.ExtData) > 0 {
		raw, err := base64.StdEncoding.DecodeString(info.ExtData)
		if err != nil {
			log.Printf("decode data fail, err:%v", err)
			return err
		}
		rc = ioutil.NopCloser(bytes.NewReader(raw))
	} else {
		reader, err := cli.DownloadFile(ctx, &models.DownloadFileRequest{
			FileId: *c.fileid,
		})
		if err != nil {
			log.Printf("get file data for read fail, err:%v", err)
			return err
		}
		rc = reader
	}
	defer rc.Close()
	data, err := ioutil.ReadAll(rc)
	if err != nil {
		log.Printf("read file data fail, err:%v", err)
		return err
	}
	s := string(data)
	fmt.Printf("%s", s)
	return nil
}
