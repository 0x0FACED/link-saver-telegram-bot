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
}

func New(token string, apiHost string) *TelegramBot {
	logger := zaplog.New()
	bot, err := bot.New(token, opts()...)
	if err != nil {
		logger.Fatal("cant create bot instance: " + err.Error())
		return nil
	}

	client, err := api.New(apiHost)
	if err != nil {
		logger.Fatal("cant create conn with api: " + err.Error())
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
		bot.WithDefaultHandler(mainHandler),
		bot.WithCallbackQueryDataHandler("button", bot.MatchTypePrefix, callbackInlineKbHandler),
	}
}

func (b *TelegramBot) registerHandlers() {
	b.bot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler)
	b.bot.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, helpHandler)
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
