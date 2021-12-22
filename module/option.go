package module

import "tgblock/bot"

type ServiceContext struct {
	Bot         *bot.TGBotService
	MaxFileSize int64
	BlockSize   int64
	SecretId    string
	SecretKey   string
	Domain      string
	Schema      string
}

type Option func(c *ServiceContext)

func WithBot(bot *bot.TGBotService) Option {
	return func(c *ServiceContext) {
		c.Bot = bot
	}
}

func WithMaxFileSize(sz int64) Option {
	return func(c *ServiceContext) {
		c.MaxFileSize = sz
	}
}

func WithBlockSize(sz int64) Option {
	return func(c *ServiceContext) {
		c.BlockSize = sz
	}
}

func WithSecret(secretid string, secretkey string) Option {
	return func(c *ServiceContext) {
		c.SecretId = secretid
		c.SecretKey = secretkey
	}
}

func WithDomain(schema string, domain string) Option {
	return func(c *ServiceContext) {
		c.Domain = domain
		c.Schema = schema
	}
}
