package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Server          string `json:"server"`
	Secretid        string `json:"secret_id"`
	Secretkey       string `json:"secret_key"`
	MaxSigAliveTime int64  `json:"max_sig_alive_time"`
	MaxFileSize     int64  `json:"max_file_size"`
	BlockSize       int64  `json:"block_size"`
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
