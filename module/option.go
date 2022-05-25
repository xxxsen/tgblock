package module

import "github.com/xxxsen/tgblock/processor"

type ServiceContext struct {
	MaxFileSize int64
	Domain      string
	Schema      string
	Processor   *processor.FileProcessor
}

type Option func(c *ServiceContext)

func WithProcessor(proc *processor.FileProcessor) Option {
	return func(c *ServiceContext) {
		c.Processor = proc
	}
}

func WithMaxFileSize(sz int64) Option {
	return func(c *ServiceContext) {
		c.MaxFileSize = sz
	}
}

func WithDomain(schema string, domain string) Option {
	return func(c *ServiceContext) {
		c.Domain = domain
		c.Schema = schema
	}
}
