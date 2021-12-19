package main

import (
	"tgblock/bot"
	"tgblock/module"

	_ "tgblock/module/meta"
	_ "tgblock/module/upload"

	flag "github.com/xxxsen/envflag"
	"github.com/xxxsen/log"
)

var listen = flag.String("listen", ":8444", "listen address")
var token = flag.String("token", "", "bot token")
var maxFileSize = flag.Int64("max_file_size", 1*1024*1024*1024, "max size per file")
var blockSize = flag.Int64("block_size", 20*1024*1024, "block size")
var chatid = flag.Int64("chatid", 0, "chatid")
var loglevel = flag.String("log_level", "trace", "log level")

func main() {
	flag.Parse()
	log.Init("", log.StringToLevel(*loglevel), 0, 0, 0, true)

	log.Infof("LISTEN:%v", *listen)
	log.Infof("TOKEN:%v", *token)
	log.Infof("MAX_FILE_SIZE:%v", *maxFileSize)
	log.Infof("BLOCK_SIZE:%v", *blockSize)
	log.Infof("CHATID:%v", *chatid)
	log.Infof("LOG_LEVEL:%v", *loglevel)

	if len(*token) == 0 || *chatid == 0 || len(*listen) == 0 || *maxFileSize == 0 {
		log.Fatal("invalid params")
	}

	tgbot, err := bot.NewBotService(
		bot.WithChatId(*chatid),
		bot.WithToken(*token),
	)
	if err != nil {
		log.Fatalf("init tgbot fail, err:%v", err)
	}

	if err := module.Init(
		module.WithBot(tgbot),
		module.WithMaxFileSize(*maxFileSize),
		module.WithBlockSize(*blockSize),
	); err != nil {
		log.Fatalf("init module fail, err:%v", err)
	}
	log.Infof("start http server...")
	if err := module.Run(*listen); err != nil {
		log.Fatalf("run server failed")
	}
}
