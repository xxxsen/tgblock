package module

import "tgblock/bot"

type ServiceContext struct {
	Bot         *bot.TGBotService
	MaxFileSize int64
	BlockSize   int64
	AccessToken string
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

func WithAccessToken(tk string) Option {
	return func(c *ServiceContext) {
		c.AccessToken = tk
	}
}
