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
	req := &gen.SaveLinkRequest{
		OriginalUrl: msg[1],
		Description: msg[2],
		Username:    update.Message.Chat.Username,
	}

	withTimeout, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	resp, err := api.SaveLink(withTimeout, req)

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
	logger := getLoggerFromCtx(ctx)
	logger.Debug("get links: "+update.Message.Text, zap.String("user", update.Message.Chat.Username))
	// TODO: logic of saving link with api
	api := getApiFromCtx(ctx)
	msg := strings.SplitN(update.Message.Text, " ", 2)
	// msg[0] -> /get command
	// msg[1] -> description

	req := &gen.GetLinksRequest{
		Username:    update.Message.From.Username,
		Description: msg[1],
	}

	withTimeout, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	resp, err := api.GetLinks(withTimeout, req)
	if err != nil {
		logger.Error("some error: "+err.Error(), zap.String("user", update.Message.Chat.Username), zap.String("mes", update.Message.Text))
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

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   "You selected the link: " + update.CallbackQuery.Data,
	})

	api := getApiFromCtx(ctx)

	reqParams := strings.Split(update.CallbackQuery.Data, ":")
	logger.Info("request GetLink params", zap.String("id", reqParams[1]), zap.String("desc", reqParams[2]))
	idReq, err := strconv.Atoi(reqParams[1])
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
		Username:    update.CallbackQuery.From.Username,
		Description: reqParams[2],
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

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   resp.GeneratedUrl,
	})
}

func inlineKbHandler(ctx context.Context, b *bot.Bot, update *models.Update, links []*gen.Link) {
	logger := getLoggerFromCtx(ctx)
	inline := &models.InlineKeyboardMarkup{}
	buttons := make([][]models.InlineKeyboardButton, 0, len(links))
	for _, link := range links {
		text := fmt.Sprintf("%s:%s", link.Description, link.OriginalUrl)
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
