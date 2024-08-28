package main

import (
	"flag"
	"log"

	"github.com/0x0FACED/link-saver-telegram-bot/config"
	"github.com/0x0FACED/link-saver-telegram-bot/internal/telegram"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic("Config didn't loaded.")
	}
	bot := telegram.New(cfg)
	if bot == nil {
		panic("Bot didn't created.")
	}
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
