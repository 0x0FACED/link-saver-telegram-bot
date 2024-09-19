package telegram

import (
	"context"

	"github.com/0x0FACED/link-saver-telegram-bot/config"
	api "github.com/0x0FACED/link-saver-telegram-bot/internal/grpc"
	"github.com/0x0FACED/link-saver-telegram-bot/internal/logger/zaplog"
	"github.com/go-telegram/bot"
)

type TelegramBot struct {
	bot            *bot.Bot
	apiClient      *api.APIClient
	logger         *zaplog.ZapLogger
	eventProcessor *EventProcessor
}

func New(cfg *config.Config) *TelegramBot {
	logger := zaplog.New()
	logger.Debug("Creating api client...")

	client, err := api.New(*cfg)
	if err != nil {
		logger.Fatal("cant create conn with api: " + err.Error())
		return nil
	}

	ep := NewEventProcessor(client, logger)
	bot, err := bot.New(cfg.Telegram.Token, opts(ep.mainHandler, ep.getCallback, ep.delCallback)...)
	if err != nil {
		logger.Fatal("Can't create bot instance: " + err.Error())
		return nil
	}

	return &TelegramBot{
		bot:            bot,
		apiClient:      client,
		logger:         logger,
		eventProcessor: ep,
	}

}

func opts(def bot.HandlerFunc, getCallback, delCallback bot.HandlerFunc) []bot.Option {
	return []bot.Option{
		bot.WithDefaultHandler(def),
		bot.WithCallbackQueryDataHandler("get", bot.MatchTypePrefix, getCallback),
		bot.WithCallbackQueryDataHandler("del", bot.MatchTypePrefix, delCallback),
	}
}

func (b *TelegramBot) registerHandlers() {
	b.bot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, b.eventProcessor.startHandler)
	b.bot.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, b.eventProcessor.helpHandler)
	b.bot.RegisterHandler(bot.HandlerTypeMessageText, "/save", bot.MatchTypeExact, b.eventProcessor.saveLinkHandlerHelper)
	b.bot.RegisterHandler(bot.HandlerTypeMessageText, "/get", bot.MatchTypeExact, b.eventProcessor.getLinksHandlerHelper)
	b.bot.RegisterHandler(bot.HandlerTypeMessageText, "/list", bot.MatchTypeExact, b.eventProcessor.getAllLinksHandler)
}

func (b *TelegramBot) Start() {
	b.registerHandlers()

	b.bot.Start(context.Background())
	b.logger.Info("Bot started...")
}
