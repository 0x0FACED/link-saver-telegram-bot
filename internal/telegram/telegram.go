package telegram

import (
	"context"

	api "github.com/0x0FACED/link-saver-telegram-bot/internal/grpc"
	"github.com/0x0FACED/link-saver-telegram-bot/internal/logger"
	"github.com/0x0FACED/link-saver-telegram-bot/internal/logger/zaplog"
	"github.com/go-telegram/bot"
)

type TelegramBot struct {
	bot       *bot.Bot
	apiClient *api.APIClient
	logger    logger.Logger
}

func New(token string, apiHost string) *TelegramBot {
	logger := zaplog.New()
	bot, err := bot.New(token, opts()...)
	if err != nil {
		logger.Fatal("cant create bot instance", err)
		return nil
	}

	client, err := api.New(apiHost)
	if err != nil {
		logger.Fatal("cant create conn with api", err)
		return nil
	}

	return &TelegramBot{
		bot:       bot,
		apiClient: client,
		logger:    logger,
	}

}

func opts() []bot.Option {
	return []bot.Option{
		bot.WithDefaultHandler(startHandler),
	}
}

func (b *TelegramBot) registerHandlers() {
	b.bot.RegisterHandler(bot.HandlerTypeMessageText, "/save", bot.MatchTypeExact, saveLinkHandler)
	b.bot.RegisterHandler(bot.HandlerTypeMessageText, "/delete", bot.MatchTypeExact, deleteLinkHandler)
	b.bot.RegisterHandler(bot.HandlerTypeMessageText, "/get", bot.MatchTypeExact, getLinksHandler)
}

func (b *TelegramBot) Start() {
	b.registerHandlers()
	b.bot.Start(b.getCtx())
}

type loggerKey string
type apiKey string

func (b *TelegramBot) getCtx() context.Context {
	ctx := context.WithValue(context.Background(), loggerKey("logger"), b.logger)
	ctx = context.WithValue(ctx, apiKey("api"), b.apiClient)

	return ctx
}
