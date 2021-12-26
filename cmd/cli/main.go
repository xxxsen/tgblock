package main

import (
	"context"
	"flag"
	"log"
	"os"
	"runtime"

	"github.com/xxxsen/tgblock/client"
	"github.com/xxxsen/tgblock/command"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("specify operation plz")
	}
	cmd := os.Args[1]
	executor, ok := command.Get(cmd)
	if !ok {
		log.Fatalf("cmd:%s not found", cmd)
	}
	fg := flag.NewFlagSet(cmd, flag.ExitOnError)
	executor.Args(fg)
	fg.Parse(os.Args[2:])

	if !executor.Check() {
		log.Fatalf("check cmd args fail")
	}

	conf := os.Getenv("CMD_TGBLOCK_CLIENT_CONF")
	if len(conf) == 0 {
		if runtime.GOOS == "windows" {
			conf = "c:\\tgblock\\client.json"
		} else {
			conf = "/etc/tgblock/client.json"
		}
	}
	c, err := ParseFile(conf)
	if err != nil {
		log.Fatalf("client config parse fail, conf:%s", conf)
	}
	cli, err := client.New(
		client.WithAddress(c.Server),
		client.WithSecret(c.Secretid, c.Secretkey),
		client.WithMaxSigAliveTime(c.MaxSigAliveTime),
		client.WithFileSize(c.MaxFileSize),
		client.WithBlockSize(c.BlockSize),
	)
	if err != nil {
		log.Fatalf("init client fail, err:%v", err)
	}
	if err := executor.Exec(context.Background(), cli); err != nil {
		log.Printf("exec cmd:%s fail, err:%v", cmd, err)
	}
}
