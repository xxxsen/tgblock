package client

type Config struct {
	Secretid        string
	Secretkey       string
	Address         string
	BlockSize       int64
	MaxFileSize     int64
	MaxSigAliveTime int64
}

type Option func(c *Config)

func WithMaxSigAliveTime(sec int64) Option {
	return func(c *Config) {
		c.MaxSigAliveTime = sec
	}
}

func WithSecret(id string, key string) Option {
	return func(c *Config) {
		c.Secretid = id
		c.Secretkey = key
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
