package telegram

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

func saveLinkHandlerHelper(ctx context.Context, b *bot.Bot, update *models.Update) {
	logger := getLoggerFromCtx(ctx)
	logger.Debug("start: "+update.Message.Text, zap.String("user", update.Message.Chat.Username))
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      saveMsg_EN,
		ParseMode: models.ParseModeMarkdownV1,
	})
}

func getLinksHandlerHelper(ctx context.Context, b *bot.Bot, update *models.Update) {
	logger := getLoggerFromCtx(ctx)
	logger.Debug("start: "+update.Message.Text, zap.String("user", update.Message.Chat.Username))
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      getMsg_EN,
		ParseMode: models.ParseModeMarkdownV1,
	})
}
