package command

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"tgblock/client"
	"tgblock/module/models"
)

func init() {
	Regist(&CmdInfo{})
}

type CmdInfo struct {
	fileid    *string
	printJson *bool
}

func (c *CmdInfo) Name() string {
	return "info"
}

func (c *CmdInfo) Args(f *flag.FlagSet) {
	c.fileid = f.String("fileid", "", "file id")
	c.printJson = f.Bool("print_json", true, "is print json")
}

func (c *CmdInfo) Check() bool {
	if len(*c.fileid) == 0 {
		log.Printf("fileid is nil")
		return false
	}
	return true
}

func (c *CmdInfo) Exec(ctx context.Context, cli *client.Client) error {
	info, err := cli.GetFileInfo(ctx, &models.GetFileInfoRequest{
		FileId: *c.fileid,
	})
	if err != nil {
		return err
	}
	js, err := json.Marshal(info)
	if err != nil {
		log.Printf("marshal response to json fail, err:%v", err)
		return err
	}
	if *c.printJson {
		fmt.Printf("%s", string(js))
		return nil
	}
	fmt.Printf("Name:%s\n", info.FileName)
	fmt.Printf("File Size:%d\n", info.FileSize)
	fmt.Printf("Hash:%s\n", info.Hash)
	fmt.Printf("Block Size:%d\n", info.BlockSize)
	fmt.Printf("Block Count:%d\n", info.BlockCount)
	fmt.Printf("Create Time:%d\n", info.CreateTime)
	fmt.Printf("Finish Time:%d\n", info.FinishTime)
	fmt.Printf("File Mode:%d\n", info.FileMode)
	fmt.Printf("===BLOCK HASH===\n")
	for index, item := range info.BlockHash {
		fmt.Printf("IDX:%d, HASH:%s\n", index, item)
	}
	return nil
}
