package client

type Config struct {
	AccessToken string
	Address     string
	BlockSize   int64
	MaxFileSize int64
}

type Option func(c *Config)

func WithAccessToken(s string) Option {
	return func(c *Config) {
		c.AccessToken = s
	}
}

func WithAddress(addr string) Option {
	return func(c *Config) {
		c.Address = addr
	}
}

func WithBlockSize(blk int64) Option {
	return func(c *Config) {
		c.BlockSize = blk
	}
}

func WithFileSize(fz int64) Option {
	return func(c *Config) {
		c.MaxFileSize = fz
	}
}
