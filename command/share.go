package command

import (
	"context"
	"flag"
	"log"
	"tgblock/client"
	"tgblock/module/models"
	"time"
)

func init() {
	Regist(&CmdShare{})
}

type CmdShare struct {
	fileid       *string
	file         *string
	expirehour   *int64
	expireSecond *int64
	expireMinute *int64
	key          *string
}

func (c *CmdShare) Name() string {
	return "share"
}

func (c *CmdShare) Args(f *flag.FlagSet) {
	c.file = f.String("file", "", "upload file and create share")
	c.fileid = f.String("fileid", "", "use fileid to create share")
	c.expirehour = f.Int64("hour", 0, "expire hour")
	c.expireSecond = f.Int64("second", 0, "expire second")
	c.expireMinute = f.Int64("minute", 0, "expire minute")
	c.key = f.String("key", "abcd", "encrypt key")
}

func (c *CmdShare) Check() bool {
	if len(*c.fileid) == 0 && len(*c.file) == 0 {
		log.Printf("both fileid/file are nil")
		return false
	}
	if len(*c.key) == 0 {
		log.Printf("key should not be nil")
		return false
	}
	return true
}

func (c *CmdShare) Exec(ctx context.Context, cli *client.Client) error {
	fileid := *c.fileid
	if len(fileid) == 0 {
		fid, err := cli.BlockUploadFile(ctx, *c.file)
		if err != nil {
			log.Printf("block upload failed, err:%v", err)
			return err
		}
		fileid = fid
	}
	sec := int64(*c.expirehour*3600 + *c.expireMinute*60 + *c.expireSecond)
	ts := time.Now().Unix() + sec
	if sec == 0 {
		ts = 0
	}
	rsp, err := cli.CreateShare(ctx, &models.CreateShareRequest{
		FileId:     fileid,
		Key:        *c.key,
		ExpireTime: ts,
	})
	if err != nil {
		log.Printf("create share link fail, err:%v", err)
		return err
	}
	log.Printf("create share succ, url:%s", rsp.URL)
	return nil
}
