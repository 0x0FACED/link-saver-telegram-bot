package telegram

import (
	"context"
	"strings"
	"time"

	"github.com/0x0FACED/link-saver-telegram-bot/internal/grpc"
	"github.com/0x0FACED/link-saver-telegram-bot/internal/logger/zaplog"
	"github.com/0x0FACED/proto-files/link_service/gen"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

var startMsg = `Hello!

This bot can:

1. Save links you send: <link> <description>
Example: https://youtube.com/example/link youtube link
2. Gave you links you saved: <description>
Example: youtube link
Bot will return you: 
https://youtube.com/example/link
<generated link in the server>
3. Delete links you saved: <description>
`

func mainHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	logger := getLoggerFromCtx(ctx)
	logger.Debug("new message: "+update.Message.Text, zap.String("user", update.Message.Chat.Username))
	msgParts := strings.Fields(update.Message.Text)
	switch msgParts[0] {
	case "/start":
		startHandler(ctx, b, update)
	case "/help":
		helpHandler(ctx, b, update)
	case "/get":
		getLinksHandler(ctx, b, update)
	case "/save":
		saveLinkHandler(ctx, b, update)
	case "/delete":
		deleteLinkHandler(ctx, b, update)
	default:
		helpHandler(ctx, b, update)
	}
}

func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	logger := getLoggerFromCtx(ctx)
	logger.Debug("start: "+update.Message.Text, zap.String("user", update.Message.Chat.Username))
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   startMsg,
	})
}

func saveLinkHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	logger := getLoggerFromCtx(ctx)
	logger.Debug("save: "+update.Message.Text, zap.String("user", update.Message.Chat.Username))
	// TODO: logic of saving link with api
	api := getApiFromCtx(ctx)
	msg := strings.SplitN(update.Message.Text, " ", 3)
	logger.Debug("part0: "+msg[0], zap.String("user", update.Message.Chat.Username))
	logger.Debug("part1: "+msg[1], zap.String("user", update.Message.Chat.Username))
	logger.Debug("part2: "+msg[2], zap.String("user", update.Message.Chat.Username))
	req := gen.SaveLinkRequest{
		OriginalUrl: msg[1],
		Description: msg[2],
		Username:    update.Message.Chat.Username,
	}

	withTimeout, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	resp, err := api.SaveLink(withTimeout, &req)

	if err != nil {
		logger.Error("some error: "+err.Error(), zap.String("user", update.Message.Chat.Username), zap.String("mes", update.Message.Text))
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Error! Try again.",
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Saved. " + resp.Message,
	})
}

func deleteLinkHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// TODO: logic of deleting link by api
}

func getLinksHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// TODO: logic of get links by api
}

func helpHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	logger := getLoggerFromCtx(ctx)
	logger.Debug("new message: "+update.Message.Text, zap.String("user", update.Message.Chat.Username))
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   helpMsg,
	})
}

func getLoggerFromCtx(ctx context.Context) *zaplog.ZapLogger {
	return ctx.Value(loggerKey("logger")).(*zaplog.ZapLogger)
}

func getApiFromCtx(ctx context.Context) *grpc.APIClient {
	return ctx.Value(apiKey("api")).(*grpc.APIClient)
}
