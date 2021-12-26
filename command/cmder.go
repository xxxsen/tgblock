package command

import (
	"context"
	"flag"
	"fmt"

	"github.com/xxxsen/tgblock/client"
)

type Command interface {
	Name() string
	Args(f *flag.FlagSet)
	Check() bool
	Exec(ctx context.Context, cli *client.Client) error
}

var cmdMap = make(map[string]Command)

func Regist(cmd Command) {
	if _, ok := cmdMap[cmd.Name()]; ok {
		panic(fmt.Errorf("cmd:%s exists", cmd.Name()))
	}
	cmdMap[cmd.Name()] = cmd
}

func Get(name string) (Command, bool) {
	c, ok := cmdMap[name]
	return c, ok
}
