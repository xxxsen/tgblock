package command

import (
	"context"
	"flag"
	"log"

	"github.com/xxxsen/tgblock/client"
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
	fileid, err := cli.BlockUploadFile(ctx, *c.file)
	if err != nil {
		log.Printf("upload file fail, err:%v", err)
		return err
	}
	log.Printf("upload succ, fileid:%s", fileid)
	return nil
}
