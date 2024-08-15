package main

import (
	"flag"
	"log"

	"github.com/0x0FACED/link-saver-telegram-bot/internal/telegram"
)

func main() {
	bot := telegram.New(mustToken(), "aasdasd")

	bot.Start()
}

func mustToken() string {
	token := flag.String("token", "", "token for telegram bot")

	flag.Parse()

	if *token == "" {
		log.Fatalln("token must be specified")
	}

	return *token
}
