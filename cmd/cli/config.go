package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Server      string `json:"server"`
	AccessToken string `json:"access_token"`
	MaxFileSize int64  `json:"max_file_size"`
	BlockSize   int64  `json:"block_size"`
}

func ParseFile(f string) (*Config, error) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	if err := json.Unmarshal(data, c); err != nil {
		return nil, err
	}
	return c, nil
}
