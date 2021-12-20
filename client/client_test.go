package client

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"testing"
	"tgblock/module/download"
	"tgblock/module/meta"
	"tgblock/module/sys"
	"tgblock/processor"

	"github.com/stretchr/testify/assert"
)

var testFileId = "ABj4C6PqYcBg9lWhAbvIAAEEpAAA-uoTMoitulXfb1fNygokdD2CcYbMAADkxAAUgACAQB"

func getClient() *Client {
	cli, err := New(WithAddress("http://127.0.0.1:8444"), WithAccessToken("abc"))
	if err != nil {
		panic(err)
	}
	return cli
}

func TestGetSysInfo(t *testing.T) {
	client := getClient()
	rsp, err := client.GetSysInfo(context.Background(), &sys.GetSysInfoRequest{})
	assert.NoError(t, err)
	t.Logf("rsp:%+v", rsp)
}

func TestUpload(t *testing.T) {
	client := getClient()
	data := []byte("hello world, this is a test")

	r := processor.NewShaReader(bytes.NewReader(data))
	io.Copy(ioutil.Discard, r)

	rsp, err := client.BlockUpload(context.Background(), &BlockUploadRequest{
		Name:   "hello.txt",
		Size:   int64(len(data)),
		Reader: bytes.NewReader(data),
		Hash:   r.GetSum(),
	})
	assert.NoError(t, err)
	t.Logf("rsp:%+v", rsp)
}

func TestDownload(t *testing.T) {
	client := getClient()
	rsp, err := client.DownloadFile(context.Background(), &download.DownloadFileRequest{
		FileId: testFileId,
	})
	assert.NoError(t, err)
	defer rsp.Close()
	data, err := ioutil.ReadAll(rsp)
	assert.NoError(t, err)
	t.Logf("data:%+v", string(data))
}

func TestDownloadBlock(t *testing.T) {
	client := getClient()
	rsp, err := client.DownloadBlock(context.Background(), &download.DownloadBlockRequest{
		FileId:     testFileId,
		BlockIndex: 0,
	})
	assert.NoError(t, err)
	defer rsp.Close()
	data, err := ioutil.ReadAll(rsp)
	assert.NoError(t, err)
	t.Logf("data:%+v", string(data))
}

func TestGetFileInfo(t *testing.T) {
	client := getClient()
	rsp, err := client.GetFileInfo(context.Background(), &meta.GetFileInfoRequest{
		FileId: testFileId,
	})
	assert.NoError(t, err)
	t.Logf("rsp:%+v", rsp)
}

func TestUploadBigBlock(t *testing.T) {
	client := getClient()
	data := make([]byte, 41*1024*1024)
	for i := 0; i < len(data); i++ {
		data[i] = byte(i) % 255
	}

	r := processor.NewShaReader(bytes.NewReader(data))
	io.Copy(ioutil.Discard, r)

	rsp, err := client.BlockUpload(context.Background(), &BlockUploadRequest{
		Name:   "bigfile.txt",
		Size:   int64(len(data)),
		Reader: bytes.NewReader(data),
		Hash:   r.GetSum(),
	})
	assert.NoError(t, err)
	t.Logf("rsp:%+v", rsp)
}
