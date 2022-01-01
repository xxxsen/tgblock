package bot

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	lru "github.com/hnlq715/golang-lru"
	"github.com/xxxsen/log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
)

type TGBotService struct {
	cfg      *TGBotConfig
	bot      *tgbotapi.BotAPI
	client   *http.Client
	cacheLnk *lru.Cache
}

func NewBotService(opts ...Option) (*TGBotService, error) {
	c := &TGBotConfig{}
	for _, opt := range opts {
		opt(c)
	}
	s := &TGBotService{cfg: c}
	return s, s.init()
}

func (s *TGBotService) asyncUpdate() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := s.bot.GetUpdatesChan(u)
	for update := range updates {
		log.Tracef("recv message from remote, sender:%d, msg:%s",
			update.Message.Chat.ID, update.Message.Text)
	}
	return nil

}

func (s *TGBotService) init() error {
	s.cacheLnk, _ = lru.New(20000)
	//init http client
	s.client = &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 10 * time.Second,
			}).Dial,
			IdleConnTimeout: 20 * time.Second,
			MaxIdleConns:    20,
		},
	}

	//parse config
	bot, err := tgbotapi.NewBotAPI(s.cfg.Token)
	if err != nil {
		return err
	}
	s.bot = bot
	go func() {
		err := s.asyncUpdate()
		if err != nil {
			log.Errorf("async update bot fail, err:%v", err)
		}
	}()
	return nil
}

func (s *TGBotService) Upload(ctx context.Context, size int64, reader io.Reader) (string, error) {
	sname := uuid.NewString()
	freader := tgbotapi.FileReader{
		Name:   sname,
		Reader: reader,
	}
	doc := tgbotapi.NewDocument(s.cfg.ChatID, freader)
	doc.DisableNotification = true
	msg, err := s.bot.Send(doc)
	if err != nil {
		return "", err
	}
	log.Tracef("upload file to robot, name:%s, size:%d, file info:%+v", sname, size, *msg.Document)
	return msg.Document.FileID, nil
}

func (s *TGBotService) cacheGetURL(ctx context.Context, hash string) (string, error) {
	if lnk, ok := s.cacheLnk.Get(hash); ok {
		return lnk.(string), nil
	}

	cf := tgbotapi.FileConfig{FileID: hash}
	f, err := s.bot.GetFile(cf)
	if err != nil {
		return "", err
	}
	lnk := f.Link(s.bot.Token)
	//这里应该能1小时有效的...
	s.cacheLnk.AddEx(hash, lnk, 30*time.Minute)
	return lnk, nil
}

func (s *TGBotService) Download(ctx context.Context, hash string) (io.ReadCloser, error) {
	return s.DownloadAt(ctx, hash, 0)
}

func (s *TGBotService) DownloadAt(ctx context.Context, hash string, index int64) (io.ReadCloser, error) {
	lnk, err := s.cacheGetURL(ctx, hash)
	if err != nil {
		return nil, err
	}
	log.Debugf("read hash:%s, link:%s", hash, lnk)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, lnk, nil)
	if err != nil {
		return nil, err
	}
	if index != 0 {
		rangeHeader := fmt.Sprintf("bytes=%d-", index)
		req.Header.Set("Range", rangeHeader)
	}
	rsp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	//caller should close rsp.Body
	if rsp.StatusCode/100 != 2 {
		rsp.Body.Close()
		return nil, fmt.Errorf("status code not ok, code:%d", rsp.StatusCode)
	}
	if index != 0 && len(rsp.Header.Get("Content-Range")) == 0 {
		rsp.Body.Close()
		return nil, fmt.Errorf("not support range")
	}

	return rsp.Body, nil
}
