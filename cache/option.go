package cache

type Config struct {
	memKeySize  int64
	fileKeySize int64
	codec       Codec
	caller      OnCacheMissFunc
	dir         string
	valuetype   interface{}
}

type Option func(c *Config)

func WithKeySize(mem int64, file int64) Option {
	return func(c *Config) {
		c.memKeySize = mem
		c.fileKeySize = file
	}
}

func WithCodec(codec Codec, valuetype interface{}) Option {
	return func(c *Config) {
		c.codec = codec
		c.valuetype = valuetype
	}
}

func WithCacheMissFunc(f OnCacheMissFunc) Option {
	return func(c *Config) {
		c.caller = f
	}
}

func WithTmpFileDir(dir string) Option {
	return func(c *Config) {
		c.dir = dir
	}
}
