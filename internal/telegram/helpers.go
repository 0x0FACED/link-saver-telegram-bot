package telegram

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"strings"
	"time"

	"github.com/0x0FACED/link-saver-telegram-bot/internal/utils"
	"github.com/0x0FACED/proto-files/link_service/gen"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

func (h *EventProcessor) saveLinkHandlerHelper(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      saveMsg_EN,
		ParseMode: models.ParseModeMarkdown,
	})
}

func (h *EventProcessor) savePDFHandlerHelper(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      savePDFMsg_EN,
		ParseMode: models.ParseModeMarkdown,
	})
}

func (h *EventProcessor) getLinksHandlerHelper(ctx context.Context, b *bot.Bot, update *models.Update) {
	h.logger.Debug("start: "+update.Message.Text, zap.String("user", update.Message.Chat.Username))
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      getMsg_EN,
		ParseMode: models.ParseModeMarkdown,
	})
}

func (h *EventProcessor) getAllLinks(ctx context.Context, b *bot.Bot, update *models.Update) []*gen.Link {
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

func decompressPDF(compressedData []byte) ([]byte, error) {
	buffer := bytes.NewBuffer(compressedData)

	reader, err := gzip.NewReader(buffer)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var decompressedData bytes.Buffer
	_, err = io.Copy(&decompressedData, reader)
	if err != nil {
		return nil, err
	}

	return decompressedData.Bytes(), nil
}

func parseMessage(text string, times int) ([]string, error) {
	msgs := strings.SplitN(text, " ", times)
	if len(msgs) != times {
		return []string{}, utils.ErrMessageFormat
	}
	return msgs, nil
}
