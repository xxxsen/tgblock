package command

import (
	"bytes"
	"context"
	"encoding/base64"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"tgblock/client"
	"tgblock/hasher"
	"tgblock/module/models"
)

func init() {
	Regist(&CmdDownload{})
}

type CmdDownload struct {
	fileid       *string
	target       *string
	createFolder *bool
}

func (c *CmdDownload) Name() string {
	return "download"
}

func (c *CmdDownload) Args(f *flag.FlagSet) {
	c.fileid = f.String("fileid", "", "file id")
	c.target = f.String("target", "./", "save location")
	c.createFolder = f.Bool("create_folder", true, "create folder if not exist")
}

func (c *CmdDownload) Check() bool {
	if len(*c.fileid) == 0 {
		log.Printf("fileid is nil")
		return false
	}
	if len(*c.target) == 0 {
		log.Printf("target is nil")
		return false
	}
	return true
}

func (c *CmdDownload) isDir(filename string) bool {
	return strings.HasSuffix(filename, "/") || strings.HasSuffix(filename, "\\")
}

func (c *CmdDownload) buildDownloadName(target string, name string) string {
	if c.isDir(target) {
		return target + string(os.PathSeparator) + name
	}
	return target
}

func (c *CmdDownload) getDir(filename string) string {
	if c.isDir(filename) {
		return filename
	}
	return filepath.Dir(filename)
}

func (c *CmdDownload) Exec(ctx context.Context, cli *client.Client) error {
	if err := os.MkdirAll(c.getDir(*c.target), 0755); err != nil {
		log.Printf("create dir fail, err:%v", err)
		return err
	}
	info, err := cli.GetFileInfo(ctx, &models.GetFileInfoRequest{
		FileId: *c.fileid,
	})
	if err != nil {
		log.Printf("read file info fail, err:%v", err)
		return err
	}
	log.Printf("read file info succ, hash:%s, size:%d, block count:%d", info.Hash, info.FileSize, info.BlockCount)
	var rc io.ReadCloser
	if len(info.ExtData) != 0 {
		raw, err := base64.StdEncoding.DecodeString(info.ExtData)
		if err != nil {
			log.Printf("decode data from extdata fail, err:%v", err)
			return err
		}
		rc = ioutil.NopCloser(bytes.NewReader(raw))
	} else {
		rc, err = cli.DownloadFile(ctx, &models.DownloadFileRequest{FileId: *c.fileid})
		if err != nil {
			log.Printf("open stream for download fail, err:%v", err)
			return err
		}
	}
	defer rc.Close()

	hashReader := hasher.NewMD5Reader(rc)
	mode := 0755
	if info.FileMode != 0 {
		mode = int(info.FileMode)
	}
	file, err := os.OpenFile(c.buildDownloadName(*c.target, info.FileName), os.O_CREATE|os.O_WRONLY, os.FileMode(mode))
	if err != nil {
		log.Printf("open file for download fail, err:%v", err)
		return err
	}
	defer file.Close()
	if _, err := io.Copy(file, hashReader); err != nil {
		log.Printf("write stream to file fail, err:%v", err)
		return err
	}
	getSum := hashReader.GetSum()
	if getSum != info.Hash {
		log.Printf("write file finish, but hash not match, expect:%s, get:%s", info.Hash, getSum)
	}
	return nil
}
