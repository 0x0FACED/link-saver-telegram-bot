package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/0x0FACED/proto-files/link_service/gen"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

func (h *EventProcessor) getCallback(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	chosen := strings.SplitN(update.CallbackQuery.Data, ":", 3)
	if len(chosen) != 3 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   "Some error, try again.",
		})
		return
	}

	h.logger.Info("request GetLink params",
		zap.String("id", chosen[1]),
		zap.String("desc", chosen[2]),
	)
	idReq, err := strconv.Atoi(chosen[1])

	if err != nil {
		h.logger.Error("error happened", zap.Error(err))

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Error! Try again." + err.Error(),
		})
		return
	}

	req := &gen.GetLinkRequest{
		UrlId:       int32(idReq),
		UserId:      update.CallbackQuery.From.ID,
		Description: chosen[2],
	}

	withTimeout, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	resp, err := h.api.GetLink(withTimeout, req)
	if err != nil {
		h.logger.Error("error happened", zap.Error(err))

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   "Error! Try again." + err.Error(),
		})
		return
	}

	text := fmt.Sprintf("*Description:* \n_%s_\n\n*Generated link:* \n%s", chosen[2], resp.GeneratedUrl)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		Text:      text,
		ParseMode: models.ParseModeMarkdownV1,
	})
}

func (h *EventProcessor) delCallback(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	chosen := strings.SplitN(update.CallbackQuery.Data, ":", 3)
	if len(chosen) != 3 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   "Some error, try again.",
		})
		return
	}

	h.logger.Info("request GetLink params",
		zap.String("id", chosen[1]),
		zap.String("desc", chosen[2]),
	)
	idReq, err := strconv.Atoi(chosen[1])

	if err != nil {
		h.logger.Error("error happened", zap.Error(err))

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Error! Try again." + err.Error(),
		})
		return
	}

	req := &gen.DeleteLinkRequest{
		LinkId: int32(idReq),
	}

	withTimeout, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	resp, err := h.api.DeleteLink(withTimeout, req)
	if err != nil {
		h.logger.Error("error happened", zap.Error(err))

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   "Error! Try again." + err.Error(),
		})
		return
	}

	text := fmt.Sprintf("*Description:* \n_%s_\n\n*Message:* `%s`", chosen[2], resp.Message)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		Text:      text,
		ParseMode: models.ParseModeMarkdownV1,
	})
}
