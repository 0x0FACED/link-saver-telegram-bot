package telegram

import (
	"context"

	api "github.com/0x0FACED/link-saver-telegram-bot/internal/grpc"
	"github.com/0x0FACED/link-saver-telegram-bot/internal/logger/zaplog"
	"github.com/go-telegram/bot"
)

type TelegramBot struct {
	bot       *bot.Bot
	apiClient *api.APIClient
	logger    *zaplog.ZapLogger
	handlers  *Handlers
}

func New(token string, apiHost string) *TelegramBot {
	logger := zaplog.New()
	logger.Debug("Creating api client...")

	client, err := api.New(apiHost)
	if err != nil {
		logger.Fatal("cant create conn with api: " + err.Error())
		return nil
	}

	handlers := NewHandlers(client, logger)
	bot, err := bot.New(token, opts(handlers.mainHandler, handlers.callbackInlineKbHandler)...)
	if err != nil {
		logger.Fatal("Can't create bot instance: " + err.Error())
		return nil
	}

	return &TelegramBot{
		bot:       bot,
		apiClient: client,
		logger:    logger,
		handlers:  handlers,
	}

}

func opts(def bot.HandlerFunc, callback bot.HandlerFunc) []bot.Option {
	return []bot.Option{
		bot.WithDefaultHandler(def),
		bot.WithCallbackQueryDataHandler("button", bot.MatchTypePrefix, callback),
	}
}

func (b *TelegramBot) registerHandlers() {
	b.bot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, b.handlers.startHandler)
	b.bot.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, b.handlers.helpHandler)
	b.bot.RegisterHandler(bot.HandlerTypeMessageText, "/save", bot.MatchTypeExact, b.handlers.saveLinkHandlerHelper)
	b.bot.RegisterHandler(bot.HandlerTypeMessageText, "/get", bot.MatchTypeExact, b.handlers.getLinksHandlerHelper)
	b.bot.RegisterHandler(bot.HandlerTypeMessageText, "/list", bot.MatchTypeExact, b.handlers.getAllLinksHandler)
}

func (b *TelegramBot) Start() {
	b.registerHandlers()

	b.bot.Start(context.Background())
}
