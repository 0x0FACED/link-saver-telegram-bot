package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/0x0FACED/link-saver-telegram-bot/internal/grpc"
	"github.com/0x0FACED/link-saver-telegram-bot/internal/logger/zaplog"
	"github.com/0x0FACED/proto-files/link_service/gen"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

func mainHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	logger := getLoggerFromCtx(ctx)
	logger.Debug("mainHandler(): "+update.Message.Text, zap.Int64("user", update.Message.From.ID))

	msgParts := strings.Fields(update.Message.Text)
	// if user sent picture or something like that, but not a text
	if len(msgParts) == 0 {
		helpHandler(ctx, b, update)
		return
	}

	switch msgParts[0] {
	case "/start":
		startHandler(ctx, b, update)
	case "/help":
		helpHandler(ctx, b, update)
	case "/get":
		getLinksHandler(ctx, b, update)
	case "/list":
		getAllLinksHandler(ctx, b, update)
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
	logger.Debug("startHandler(): "+update.Message.Text, zap.Int64("user", update.Message.From.ID))

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      startMsg_EN,
		ParseMode: models.ParseModeMarkdownV1,
	})
}

func saveLinkHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	logger := getLoggerFromCtx(ctx)
	logger.Debug("saveLinkHandler(): "+update.Message.Text, zap.Int64("user", update.Message.From.ID))

	api := getApiFromCtx(ctx)

	msg := strings.SplitN(update.Message.Text, " ", 3)
	if len(msg) != 3 {
		saveLinkHandlerHelper(ctx, b, update)
		return
	}

	req := &gen.SaveLinkRequest{
		OriginalUrl: msg[1],
		Description: msg[2],
		UserId:      update.Message.From.ID,
	}

	withTimeout, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	resp, err := api.SaveLink(withTimeout, req)

	if err != nil {
		logger.Error("error by API: "+err.Error(),
			zap.String("user", update.Message.Chat.Username),
			zap.String("user_msg", update.Message.Text),
		)
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
	logger := getLoggerFromCtx(ctx)
	logger.Debug("getLinksHandler(): "+update.Message.Text, zap.Int64("user", update.Message.From.ID))

	api := getApiFromCtx(ctx)

	msg := strings.SplitN(update.Message.Text, " ", 2)
	// msg[0] -> /get command
	// msg[1] -> description

	if len(msg) != 2 {
		getLinksHandlerHelper(ctx, b, update)
		return
	}

	req := &gen.GetLinksRequest{
		UserId:      update.Message.From.ID,
		Description: msg[1],
	}

	withTimeout, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	resp, err := api.GetLinks(withTimeout, req)

	if err != nil {
		logger.Error("some error: "+err.Error(), zap.Int64("user", update.Message.From.ID), zap.String("mes", update.Message.Text))

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Error! Try again." + err.Error(),
		})
		return
	}

	logger.Debug("got links: ", zap.Any("links", resp.Links))

	inlineKbHandler(ctx, b, update, resp.Links)
}

func getAllLinksHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	logger := getLoggerFromCtx(ctx)
	logger.Debug("getAllLinksHandler(): "+update.Message.Text, zap.Int64("user", update.Message.From.ID))

	api := getApiFromCtx(ctx)

	req := &gen.GetAllLinksRequest{
		UserId: update.Message.From.ID,
	}

	withTimeout, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	resp, err := api.GetAllLinks(withTimeout, req)
	if err != nil {
		logger.Error("some error: "+err.Error(), zap.Int64("user", update.Message.From.ID), zap.String("mes", update.Message.Text))

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Error! Try again." + err.Error(),
		})
		return
	}

	logger.Debug("got links: ", zap.Any("links", resp.Links))

	inlineKbHandler(ctx, b, update, resp.Links)
}

func callbackInlineKbHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	logger := getLoggerFromCtx(ctx)
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

	api := getApiFromCtx(ctx)

	logger.Info("request GetLink params", zap.String("id", chosen[1]), zap.String("desc", chosen[2]))
	idReq, err := strconv.Atoi(chosen[1])

	if err != nil {
		logger.Error("error happened", zap.Error(err))

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

	resp, err := api.GetLink(withTimeout, req)
	if err != nil {
		logger.Error("error happened", zap.Error(err))

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   "Error! Try again." + err.Error(),
		})
		return
	}

	text := fmt.Sprintf("*Description:* \n_%s_\n\n*Generated link:* `%s`", chosen[2], resp.GeneratedUrl)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		Text:      text,
		ParseMode: models.ParseModeMarkdownV1,
	})
}

func inlineKbHandler(ctx context.Context, b *bot.Bot, update *models.Update, links []*gen.Link) {
	if len(links) == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "No saved links 😞",
		})
	}

	logger := getLoggerFromCtx(ctx)

	inline := &models.InlineKeyboardMarkup{}
	buttons := make([][]models.InlineKeyboardButton, 0, len(links))
	for i, link := range links {
		text := fmt.Sprintf("%d. %s: %s", i, link.Description, link.OriginalUrl)
		callbackData := fmt.Sprintf("button:%d:%s", link.LinkId, link.Description)

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
		logger.Error("err while send message (keyboard)", zap.Error(err))
	}
}

func helpHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	logger := getLoggerFromCtx(ctx)
	logger.Debug("new message: "+update.Message.Text, zap.Int64("user", update.Message.From.ID))

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   startMsg_EN,
	})
}

func getLoggerFromCtx(ctx context.Context) *zaplog.ZapLogger {
	return ctx.Value(loggerKey("logger")).(*zaplog.ZapLogger)
}

func getApiFromCtx(ctx context.Context) *grpc.APIClient {
	return ctx.Value(apiKey("api")).(*grpc.APIClient)
}
