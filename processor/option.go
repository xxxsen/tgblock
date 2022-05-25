package processor

import (
	"github.com/xxxsen/tgblock/bot"
	"github.com/xxxsen/tgblock/cache"
	"github.com/xxxsen/tgblock/locker"
)

type Config struct {
	tgbot  *bot.TGBotService
	lcker  locker.Locker
	fcache cache.Cache
	seckey string
}

type Option func(c *Config)

func WithBot(tgbot *bot.TGBotService) Option {
	return func(c *Config) {
		c.tgbot = tgbot
	}
}

func WithLocker(lck locker.Locker) Option {
	return func(c *Config) {
		c.lcker = lck
	}
}

func WithCache(fc cache.Cache) Option {
	return func(c *Config) {
		c.fcache = fc
	}
}

func WithSecKey(key string) Option {
	return func(c *Config) {
		c.seckey = key
	}
}
