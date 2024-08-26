package telegram

import (
	"context"
	"time"

	"github.com/0x0FACED/proto-files/link_service/gen"
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

func (h *Handlers) getAllLinks(ctx context.Context, b *bot.Bot, update *models.Update) []*gen.Link {
	req := &gen.GetAllLinksRequest{
		UserId: update.Message.From.ID,
	}

	withTimeout, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	resp, err := h.api.GetAllLinks(withTimeout, req)
	if err != nil {
		h.logger.Error("some error: "+err.Error(), zap.Int64("user", update.Message.From.ID), zap.String("mes", update.Message.Text))

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Error! Try again." + err.Error(),
		})
		return nil
	}

	h.logger.Debug("got links: ", zap.Any("links", resp.Links))

	return resp.Links
}
