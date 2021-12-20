package client

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"testing"
	"tgblock/module/sys"
	"tgblock/processor"

	"github.com/stretchr/testify/assert"
)

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
