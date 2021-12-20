package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	codec "tgblock/coder/client"
	"tgblock/module/download"
	"tgblock/module/meta"
	"tgblock/module/sys"
	"tgblock/module/upload"
	"tgblock/processor"
	"time"
)

type Client struct {
	c      *Config
	client *http.Client
}

func New(opts ...Option) (*Client, error) {
	c := &Config{}
	for _, opt := range opts {
		opt(c)
	}
	if strings.HasSuffix(c.Address, "/") {
		c.Address = strings.TrimRight(c.Address, "/")
	}
	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 10 * time.Second,
			}).Dial,
			IdleConnTimeout: 20 * time.Second,
			MaxIdleConns:    5,
		},
	}

	cli := &Client{c: c, client: client}
	if c.BlockSize == 0 || c.MaxFileSize == 0 {
		info, err := cli.GetSysInfo(context.Background(), &sys.GetSysInfoRequest{})
		if err != nil {
			return nil, err
		}
		c.BlockSize = info.BlockSize
		c.MaxFileSize = info.MaxFileSize
	}
	return cli, nil
}

func (c *Client) buildURL(api string) string {
	return c.c.Address + api
}

func (c *Client) buildRequest(method string, api string, codec codec.Encoder, input interface{}) (*http.Request, error) {
	api = c.buildURL(api)
	req, err := http.NewRequest(method, api, nil)
	if err != nil {
		return nil, err
	}
	if err := codec.Encode(req, input); err != nil {
		return nil, err
	}
	req.Header.Set("acess_token", c.c.AccessToken)
	return req, nil
}

func (c *Client) callRequest(req *http.Request, codec codec.Decoder, response interface{}) error {
	rsp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("status not ok, code:%d", rsp.StatusCode)
	}
	return codec.Decode(rsp, response)
}

func (c *Client) call(method, api string, codec codec.Codec, input interface{}, output interface{}) error {
	req, err := c.buildRequest(method, api, codec, input)
	if err != nil {
		return fmt.Errorf("build request fail, err:%v", err)
	}
	if err := c.callRequest(req, codec, output); err != nil {
		return fmt.Errorf("build response fail, err:%v", err)
	}
	return nil
}

func (c *Client) GetFileInfo(ctx context.Context, request *meta.GetFileInfoRequest) (*meta.GetFileInfoResponse, error) {

	rsp := &meta.GetFileInfoResponse{}
	codec := codec.MakeCodec(codec.DefaultURLCodec, codec.DefaultJsonCodec)

	if err := c.call(http.MethodGet, apiGetFileMeta, codec, request, rsp); err != nil {
		return nil, err
	}
	return rsp, nil
}

func (c *Client) DownloadBlock(ctx context.Context, request *download.DownloadBlockRequest) (io.ReadCloser, error) {
	var rc io.ReadCloser
	codec := codec.MakeCodec(codec.DefaultURLCodec, codec.DefaultStreamCodec)

	if err := c.call(http.MethodGet, apiDownloadBlock, codec, request, &rc); err != nil {
		return nil, err
	}
	return rc, nil
}

func (c *Client) DownloadFile(ctx context.Context, request *download.DownloadFileRequest) (io.ReadCloser, error) {
	var rc io.ReadCloser

	codec := codec.MakeCodec(codec.DefaultURLCodec, codec.DefaultStreamCodec)
	if err := c.call(http.MethodGet, apiDownloadFile, codec, request, &rc); err != nil {
		return nil, err
	}
	return rc, nil
}

func (c *Client) BlockUploadBegin(ctx context.Context, request *upload.BlockUploadBeginRequest) (*upload.BlockUploadBeginResponse, error) {
	rsp := &upload.BlockUploadBeginResponse{}

	if err := c.call(http.MethodPost, apiBlockUploadBegin, codec.DefaultJsonCodec, request, rsp); err != nil {
		return nil, err
	}
	return rsp, nil
}

func (c *Client) BlockUploadPart(ctx context.Context, uploadid string, file *FileInfo) error {
	m := make(map[string]interface{})
	m["uploadid"] = uploadid
	m["hash"] = file.Hash
	m["file"] = &codec.FormFileInfo{
		Name: file.Name,
		Size: file.Size,
		File: file.File,
	}
	if err := c.call(http.MethodPost, apiBlockUploadPart, codec.DefaultFormFileCodec, m, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) BlockUploadEnd(ctx context.Context, request *upload.BlockUploadEndRequest) (*upload.BlockUploadEndResponse, error) {
	rsp := &upload.BlockUploadEndResponse{}
	if err := c.call(http.MethodPost, apiBlockUploadEnd, codec.DefaultJsonCodec, request, rsp); err != nil {
		return nil, err
	}
	return rsp, nil
}

func (c *Client) BlockUpload(ctx context.Context, request *BlockUploadRequest) (*BlockUploadResponse, error) {
	if request.Size == 0 {
		return nil, fmt.Errorf("empty file")
	}
	begin, err := c.BlockUploadBegin(ctx, &upload.BlockUploadBeginRequest{
		Name:     request.Name,
		FileSize: request.Size,
		Hash:     request.Hash,
	})
	if err != nil {
		return nil, fmt.Errorf("block upload begin fail, err:%v", err)
	}
	maxBlock := (request.Size + c.c.BlockSize - 1) / c.c.BlockSize
	for i := 0; i < int(maxBlock); i++ {
		size := c.c.BlockSize
		if i == int(maxBlock)-1 {
			size = request.Size - request.Size/c.c.BlockSize*c.c.BlockSize
		}
		limitReader := processor.NewShaReader(io.LimitReader(request.Reader, size))
		data, err := ioutil.ReadAll(limitReader)
		if err != nil {
			return nil, fmt.Errorf("read part fail, err:%v", err)
		}
		if err := c.BlockUploadPart(ctx, begin.UploadId, &FileInfo{
			Name: fmt.Sprintf("%s.part.%d", request.Name, i),
			Size: size,
			Hash: limitReader.GetSum(),
			File: bytes.NewReader(data),
		}); err != nil {
			return nil, fmt.Errorf("upload part fail, err:%v", err)
		}
	}
	end, err := c.BlockUploadEnd(ctx, &upload.BlockUploadEndRequest{
		UploadId: begin.UploadId,
	})
	if err != nil {
		return nil, fmt.Errorf("end upload fail, err:%v", err)
	}
	return &BlockUploadResponse{
		FileId: end.FileId,
	}, nil
}

func (c *Client) GetSysInfo(ctx context.Context, request *sys.GetSysInfoRequest) (*sys.GetSysInfoResponse, error) {
	if request.Timestamp == 0 {
		request.Timestamp = time.Now().Unix()
	}
	rsp := &sys.GetSysInfoResponse{}
	if err := c.call(http.MethodGet, apiGetSysInfo, codec.MakeCodec(codec.DefaultURLCodec, codec.DefaultJsonCodec), request, rsp); err != nil {
		return nil, err
	}
	return rsp, nil
}