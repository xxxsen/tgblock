package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	codec "github.com/xxxsen/tgblock/coder/client"
	"github.com/xxxsen/tgblock/hasher"
	"github.com/xxxsen/tgblock/module/models"
	"github.com/xxxsen/tgblock/security"
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
		info, err := cli.GetSysInfo(context.Background(), &models.GetSysInfoRequest{})
		if err != nil {
			return nil, err
		}
		c.BlockSize = info.BlockSize
		c.MaxFileSize = info.MaxFileSize
	}
	return cli, nil
}

func (c *Client) MaxFileSize() int64 {
	return c.c.MaxFileSize
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
	sec := c.c.MaxSigAliveTime
	if sec == 0 {
		sec = 60
	}
	timestamp := time.Now().Unix() + sec

	sig, err := security.CreateSig(c.c.Secretid, c.c.Secretkey, timestamp)
	if err != nil {
		return nil, err
	}
	req.Header.Set(security.SigSecretId, c.c.Secretid)
	req.Header.Set(security.SigSecretTs, fmt.Sprintf("%d", timestamp))
	req.Header.Set(security.SigSecretSig, sig)
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

func (c *Client) GetFileInfo(ctx context.Context, request *models.GetFileInfoRequest) (*models.GetFileInfoResponse, error) {

	rsp := &models.GetFileInfoResponse{}
	codec := codec.MakeCodec(codec.DefaultURLCodec, codec.DefaultJsonCodec)

	if err := c.call(http.MethodGet, apiGetFileMeta, codec, request, rsp); err != nil {
		return nil, err
	}
	return rsp, nil
}

func (c *Client) DownloadFile(ctx context.Context, request *models.DownloadFileRequest) (io.ReadCloser, error) {
	var rc io.ReadCloser

	codec := codec.MakeCodec(codec.DefaultURLCodec, codec.DefaultStreamCodec)
	if err := c.call(http.MethodGet, apiDownloadFile, codec, request, &rc); err != nil {
		return nil, err
	}
	return rc, nil
}

func (c *Client) BlockUploadBegin(ctx context.Context, request *models.BlockUploadBeginRequest) (*models.BlockUploadBeginResponse, error) {
	rsp := &models.BlockUploadBeginResponse{}

	if err := c.call(http.MethodPost, apiBlockUploadBegin, codec.DefaultJsonCodec, request, rsp); err != nil {
		return nil, err
	}
	return rsp, nil
}

func (c *Client) BlockUploadPart(ctx context.Context, uploadid string, index int64, file *FileInfo) error {
	m := make(map[string]interface{})
	m["uploadid"] = uploadid
	m["hash"] = file.Hash
	m["index"] = index
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

func (c *Client) BlockUploadEnd(ctx context.Context, request *models.BlockUploadEndRequest) (*models.BlockUploadEndResponse, error) {
	rsp := &models.BlockUploadEndResponse{}
	if err := c.call(http.MethodPost, apiBlockUploadEnd, codec.DefaultJsonCodec, request, rsp); err != nil {
		return nil, err
	}
	return rsp, nil
}

func (c *Client) BlockUploadFile(ctx context.Context, file string) (string, error) {
	stat, err := os.Stat(file)
	if err != nil {
		return "", err
	}
	if stat.Size() > c.MaxFileSize() {
		return "", fmt.Errorf("size exceed, max:%d", c.MaxFileSize())
	}
	md5str, err := hasher.CalcMD5(file)
	if err != nil {
		return "", err
	}
	reader, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer reader.Close()
	req := &BlockUploadRequest{
		Name:   filepath.Base(file),
		Hash:   md5str,
		Size:   stat.Size(),
		Reader: reader,
		Mode:   int64(stat.Mode()),
	}
	if stat.Size() == 0 {
		req.ForceZero = true
	}
	rsp, err := c.BlockUpload(ctx, req)
	if err != nil {
		return "", err
	}
	return rsp.FileId, nil
}

func (c *Client) BlockUpload(ctx context.Context, request *BlockUploadRequest) (*BlockUploadResponse, error) {
	if request.Size == 0 && !request.ForceZero {
		return nil, fmt.Errorf("empty file")
	}
	begin, err := c.BlockUploadBegin(ctx, &models.BlockUploadBeginRequest{
		Name:      request.Name,
		FileSize:  request.Size,
		Hash:      request.Hash,
		FileMode:  request.Mode,
		ForceZero: request.ForceZero,
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
		limitReader := hasher.NewMD5Reader(io.LimitReader(request.Reader, size))
		data, err := ioutil.ReadAll(limitReader)
		if err != nil {
			return nil, fmt.Errorf("read part fail, err:%v", err)
		}
		if err := c.BlockUploadPart(ctx, begin.UploadId, int64(i), &FileInfo{
			Name: fmt.Sprintf("%s.part.%d", request.Name, i),
			Size: size,
			Hash: limitReader.GetSum(),
			File: bytes.NewReader(data),
		}); err != nil {
			return nil, fmt.Errorf("upload part fail, err:%v", err)
		}
	}
	end, err := c.BlockUploadEnd(ctx, &models.BlockUploadEndRequest{
		UploadId: begin.UploadId,
	})
	if err != nil {
		return nil, fmt.Errorf("end upload fail, err:%v", err)
	}
	return &BlockUploadResponse{
		FileId: end.FileId,
	}, nil
}

func (c *Client) GetSysInfo(ctx context.Context, request *models.GetSysInfoRequest) (*models.GetSysInfoResponse, error) {
	if request.Timestamp == 0 {
		request.Timestamp = time.Now().Unix()
	}
	rsp := &models.GetSysInfoResponse{}
	if err := c.call(http.MethodGet, apiGetSysInfo, codec.MakeCodec(codec.DefaultURLCodec, codec.DefaultJsonCodec), request, rsp); err != nil {
		return nil, err
	}
	return rsp, nil
}

func (c *Client) CreateShare(ctx context.Context, request *models.CreateShareRequest) (*models.CreateShareResponse, error) {
	if len(request.FileId) == 0 || len(request.Key) == 0 {
		return nil, fmt.Errorf("invalid fileid/key")
	}
	rsp := &models.CreateShareResponse{}
	if err := c.call(http.MethodPost, apiCreateShare, codec.DefaultJsonCodec, request, rsp); err != nil {
		return nil, err
	}
	return rsp, nil
}
