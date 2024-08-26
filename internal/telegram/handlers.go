package telegram

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/0x0FACED/link-saver-telegram-bot/internal/grpc"
	"github.com/0x0FACED/link-saver-telegram-bot/internal/logger/zaplog"
	"github.com/0x0FACED/link-saver-telegram-bot/internal/utils"
	"github.com/0x0FACED/proto-files/link_service/gen"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

type CommandHandlerFunc func(ctx context.Context, b *bot.Bot, update *models.Update)

type Handlers struct {
	handlers map[string]CommandHandlerFunc
	api      *grpc.APIClient
	logger   *zaplog.ZapLogger
}

func NewHandlers(api *grpc.APIClient, logger *zaplog.ZapLogger) *Handlers {
	h := &Handlers{
		api:      api,
		logger:   logger,
		handlers: make(map[string]CommandHandlerFunc),
	}

	h.initHandlers()

	return h
}

func (h *Handlers) initHandlers() {
	h.handlers["/start"] = h.startHandler
	h.handlers["/help"] = h.helpHandler
	h.handlers["/get"] = h.getLinksHandler
	h.handlers["/list"] = h.getAllLinksHandler
	h.handlers["/save"] = h.saveLinkHandler
	h.handlers["/delete"] = h.deleteLinkHandler
	h.handlers["/del"] = h.getAllLinksHandlerDelete
}

func (h *Handlers) mainHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update == nil || update.Message == nil {
		h.logger.Debug("mainHandler(): ", zap.String("err", "nil!"))
		return
	}
	h.logger.Debug("mainHandler(): "+update.Message.Text, zap.Int64("user", update.Message.From.ID))

	msgParts := strings.Fields(update.Message.Text)
	if len(msgParts) == 0 {
		h.helpHandler(ctx, b, update)
		return
	}

	if handler, ok := h.handlers[msgParts[0]]; ok {
		handler(ctx, b, update)
	} else {
		h.helpHandler(ctx, b, update)
	}
}

func (h *Handlers) startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update == nil || update.Message == nil || h.logger == nil {
		h.logger.Debug("mainHandler(): ", zap.String("err", "nil!"))
		return
	}
	h.logger.Debug("startHandler(): "+update.Message.Text, zap.Int64("user", update.Message.From.ID))

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      startMsg_EN,
		ParseMode: models.ParseModeMarkdown,
	})
}

func (h *Handlers) saveLinkHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	h.logger.Debug("saveLinkHandler(): "+update.Message.Text, zap.Int64("user", update.Message.From.ID))

	msg := strings.SplitN(update.Message.Text, " ", 3)
	if len(msg) != 3 {
		h.saveLinkHandlerHelper(ctx, b, update)
		return
	}

	req := &gen.SaveLinkRequest{
		OriginalUrl: msg[1],
		Description: msg[2],
		UserId:      update.Message.From.ID,
	}

	withTimeout, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	h.logger.Debug("saveLinkHandler() request: ", zap.Any("req", req))
	_, err := h.api.SaveLink(withTimeout, req)

	if err != nil {
		msg := utils.GetCodeMsgFromError(err)
		h.logger.Error("error by API",
			zap.String("user", update.Message.Chat.Username),
			zap.String("user_msg", update.Message.Text),
			zap.Error(err),
		)

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   msg + " 😞",
		})

		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Successfully Saved",
	})
}

func (h *Handlers) deleteLinkHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	h.logger.Debug("deleteLinkHandler(): "+update.Message.Text, zap.Int64("user", update.Message.From.ID))

	h.getAllLinksHandler(ctx, b, update)
}

func (h *Handlers) getLinksHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	h.logger.Debug("getLinksHandler(): "+update.Message.Text, zap.Int64("user", update.Message.From.ID))

	msg := strings.SplitN(update.Message.Text, " ", 2)
	// msg[0] -> /get command
	// msg[1] -> description

	if len(msg) != 2 {
		h.getLinksHandlerHelper(ctx, b, update)
		return
	}

	req := &gen.GetLinksRequest{
		UserId:      update.Message.From.ID,
		Description: msg[1],
	}

	withTimeout, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	resp, err := h.api.GetLinks(withTimeout, req)

	if err != nil {
		h.logger.Error("some error: "+err.Error(), zap.Int64("user", update.Message.From.ID), zap.String("mes", update.Message.Text))

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Error! Try again." + err.Error(),
		})
		return
	}

	h.logger.Debug("got links: ", zap.Any("links", resp.Links))

	h.inlineKbHandler(ctx, b, update, resp.Links, "get")
}

func (h *Handlers) getAllLinksHandlerDelete(ctx context.Context, b *bot.Bot, update *models.Update) {

	h.logger.Debug("getAllLinksHandler(): "+update.Message.Text, zap.Int64("user", update.Message.From.ID))
	links := h.getAllLinks(ctx, b, update)
	h.inlineKbHandler(ctx, b, update, links, "del")
}

func (h *Handlers) getAllLinksHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	h.logger.Debug("getAllLinksHandler(): "+update.Message.Text, zap.Int64("user", update.Message.From.ID))
	links := h.getAllLinks(ctx, b, update)
	h.inlineKbHandler(ctx, b, update, links, "get")
}

func (h *Handlers) inlineKbHandler(ctx context.Context, b *bot.Bot, update *models.Update, links []*gen.Link, tag string) {
	if len(links) == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "No saved links 😞",
		})

		return
	}

	inline := &models.InlineKeyboardMarkup{}
	buttons := make([][]models.InlineKeyboardButton, 0, len(links))
	for i, link := range links {
		text := fmt.Sprintf("%d. %s: %s", i, link.Description, link.OriginalUrl)
		callbackData := fmt.Sprintf("%s:%d:%s", tag, link.LinkId, link.Description)

		button := models.InlineKeyboardButton{
			Text:         text,
			CallbackData: callbackData,
		}

		buttons = append(buttons, []models.InlineKeyboardButton{button})
	}

	inline.InlineKeyboard = buttons

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Choose the link you want",
		ReplyMarkup: inline,
	})
	if err != nil {
		h.logger.Error("err while send message (keyboard)", zap.Error(err))
	}
}

func (h *Handlers) helpHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	h.logger.Debug("new message: "+update.Message.Text, zap.Int64("user", update.Message.From.ID))

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      helpMsg_EN,
		ParseMode: models.ParseModeMarkdown,
	})
}
