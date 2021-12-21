package command

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"tgblock/client"
	"tgblock/module/models"
	"unicode"
)

func init() {
	Regist(&CmdCat{})
}

type CmdCat struct {
	fileid        *string
	maxSizeLimit  *int64
	preventBinary *bool
}

func (c *CmdCat) Name() string {
	return "cat"
}

func (c *CmdCat) Args(f *flag.FlagSet) {
	c.fileid = f.String("fileid", "", "file id")
	c.preventBinary = f.Bool("prevent_binary", true, "is prevent binary")
	c.maxSizeLimit = f.Int64("max_size_limit", 4*1024, "max size to output")
}

func (c *CmdCat) Check() bool {
	if len(*c.fileid) == 0 {
		log.Printf("fileid is nil")
		return false
	}
	return true
}

func (c *CmdCat) isAsciiPrintable(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII || !unicode.IsPrint(r) {
			return false
		}
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
	reader, err := cli.DownloadFile(ctx, &models.DownloadFileRequest{
		FileId: *c.fileid,
	})
	if err != nil {
		log.Printf("get file data for read fail, err:%v", err)
		return err
	}
	defer reader.Close()

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Printf("read file data fail, err:%v", err)
		return err
	}
	s := string(data)
	if *c.preventBinary && !c.isAsciiPrintable(s) {
		log.Printf("binary data detect, skip")
		return nil
	}
	fmt.Printf("%s", s)
	return nil
}
