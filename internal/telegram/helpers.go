package telegram

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

func (h *Handlers) saveLinkHandlerHelper(ctx context.Context, b *bot.Bot, update *models.Update) {
	h.logger.Debug("start: "+update.Message.Text, zap.String("user", update.Message.Chat.Username))
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      saveMsg_EN,
		ParseMode: models.ParseModeMarkdown,
	})
}

func (h *Handlers) getLinksHandlerHelper(ctx context.Context, b *bot.Bot, update *models.Update) {
	h.logger.Debug("start: "+update.Message.Text, zap.String("user", update.Message.Chat.Username))
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      getMsg_EN,
		ParseMode: models.ParseModeMarkdown,
	})
}
