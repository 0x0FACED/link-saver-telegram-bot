package telegram

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/0x0FACED/link-saver-telegram-bot/internal/grpc"
	"github.com/0x0FACED/link-saver-telegram-bot/internal/logger/zaplog"
	"github.com/0x0FACED/link-saver-telegram-bot/internal/utils"
	pdf "github.com/0x0FACED/pdf-proto/pdf_service/gen"
	"github.com/0x0FACED/proto-files/link_service/gen"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

type CommandHandlerFunc func(ctx context.Context, b *bot.Bot, update *models.Update)

type EventProcessor struct {
	handlers map[string]CommandHandlerFunc
	api      *grpc.APIClient
	logger   *zaplog.ZapLogger
}

func NewEventProcessor(api *grpc.APIClient, logger *zaplog.ZapLogger) *EventProcessor {
	h := &EventProcessor{
		api:      api,
		logger:   logger,
		handlers: make(map[string]CommandHandlerFunc),
	}

	h.initHandlers()

	return h
}

func (h *EventProcessor) initHandlers() {
	h.handlers["/start"] = h.startHandler
	h.handlers["/help"] = h.helpHandler

	//h.handlers["/get"] = h.getLinksHandler
	//h.handlers["/list"] = h.getAllLinksHandler
	//h.handlers["/save"] = h.saveLinkHandler
	//h.handlers["/delete"] = h.deleteLinkHandler
	//h.handlers["/del"] = h.getAllLinksHandlerDelete

	h.handlers["/pdf"] = h.savePDFHandlerHelper
	h.handlers["/savepdf"] = h.savePDFHandler
	//h.handlers["/getpdf"] = h.getPDFHandler
	//h.handlers["/delpdf"] = h.deletePDFHandler
	//h.handlers["/flushpdfs"] = h.deleteAllPDFSHandler
}

func (h *EventProcessor) mainHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
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

func (h *EventProcessor) startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update == nil || update.Message == nil || h.logger == nil {
		h.logger.Debug("mainHandler(): ", zap.String("err", "nil!"))
		return
	}
	h.logger.Debug("startHandler(): "+update.Message.Text, zap.Int64("user", update.Message.From.ID))

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      startMsg_RU,
		ParseMode: models.ParseModeMarkdown,
	})
}

func (h *EventProcessor) saveLinkHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	h.logger.Debug("saveLinkHandler(): "+update.Message.Text, zap.Int64("user", update.Message.From.ID))

	if err := utils.ValidateSaveMessage(update.Message.Text); err != nil {
		msg := "Неверный формат сообщения. Наждите на /save, чтобы узнать, как работает команда"

		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   msg,
		})
		return
	}
	msg := strings.SplitN(update.Message.Text, " ", 3)

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
		Text:   "Успешно сохранено",
	})
}

func (h *EventProcessor) deleteLinkHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	h.logger.Debug("deleteLinkHandler(): "+update.Message.Text, zap.Int64("user", update.Message.From.ID))

	h.getAllLinksHandler(ctx, b, update)
}

func (h *EventProcessor) getLinksHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	h.logger.Debug("getLinksHandler(): "+update.Message.Text, zap.Int64("user", update.Message.From.ID))

	msg, err := parseMessage(update.Message.Text, 2)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Неверный формат сообщения",
		})
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
		h.logger.Error("some error",
			zap.Int64("user", update.Message.From.ID),
			zap.String("mes", update.Message.Text),
			zap.Error(err),
		)

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Ошибка, попробуйте позже.",
		})
		return
	}

	h.logger.Debug("got links: ", zap.Any("links", resp.Links))

	h.inlineKbHandler(ctx, b, update, resp.Links, "get")
}

func (h *EventProcessor) getAllLinksHandlerDelete(ctx context.Context, b *bot.Bot, update *models.Update) {

	h.logger.Debug("getAllLinksHandlerDelete(): "+update.Message.Text, zap.Int64("user", update.Message.From.ID))
	links := h.getAllLinks(ctx, b, update)
	h.inlineKbHandler(ctx, b, update, links, "del")
}

func (h *EventProcessor) getAllLinksHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	h.logger.Debug("getAllLinksHandler(): "+update.Message.Text, zap.Int64("user", update.Message.From.ID))
	links := h.getAllLinks(ctx, b, update)
	h.inlineKbHandler(ctx, b, update, links, "get")
}

func (h *EventProcessor) savePDFHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	h.logger.Debug("savePDFHandler(): "+update.Message.Text, zap.Int64("user", update.Message.From.ID))

	msg, err := parseMessage(update.Message.Text, 2)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Неверный формат сообщения",
		})
		return
	}

	url, err := ValidateAndFixURL(msg[1])
	if err != nil {
		if errors.Is(err, ErrInvalidURL) {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Неверный формат ссылки. Введите /pdf для получения справки.",
				ReplyParameters: &models.ReplyParameters{
					ChatID:    update.Message.Chat.ID,
					MessageID: update.Message.ID,
				},
			})
		}
		if errors.Is(err, ErrInvalidProtocol) {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Поддерживаются только 'http' и 'https' протоколы.",
				ReplyParameters: &models.ReplyParameters{
					ChatID:    update.Message.Chat.ID,
					MessageID: update.Message.ID,
				},
			})
		}
	}
	// Должен быть формат: /savepdf <link> <description> <scale> (description - название) (ПОТОМ)
	// пока что без scale и description
	// длина описания не должна быть больше 16 символов (ну условно, на первое время)
	req := &pdf.ConvertToPDFRequest{
		UserId:      update.Message.From.ID,
		OriginalUrl: url,
		Description: "tmp", // на время, пока не добавлю хранение в редисе у сервиса
		Scale:       0.7,   // на время
	}

	// Отправим сообщение, мол этот процесс может занять длительное время, ожидайте
	waitMsg, _ := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Этот процесс может занять какое-то время. Пожалуйста, подождите🙏🥺",
		ReplyParameters: &models.ReplyParameters{
			ChatID:    update.Message.Chat.ID,
			MessageID: update.Message.ID,
		},
	})

	resp, err := h.api.ConvertToPDF(ctx, req)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Ошибка, попробуйте позже.",
			ReplyParameters: &models.ReplyParameters{
				ChatID:    update.Message.Chat.ID,
				MessageID: update.Message.ID,
			},
		})
		return
	}

	decompressed, err := decompressPDF(resp.PdfData)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Ошибка, попробуйте позже.",
		})
		return
	}

	// Удаляем сообщение про ожидание
	//b.DeleteMessage(ctx, &bot.DeleteMessageParams{
	//	ChatID:    update.Message.Chat.ID,
	//	MessageID: waitMsg.ID,
	//})

	waitMsg, _ = b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    update.Message.Chat.ID,
		MessageID: waitMsg.ID,
		Text:      "Файл готов, отправляем 📥",
	})

	filepath := fmt.Sprintf("attach://%s", resp.Filename)
	_, err = b.EditMessageMedia(ctx, &bot.EditMessageMediaParams{
		ChatID:    update.Message.Chat.ID,
		MessageID: waitMsg.ID,
		Media: &models.InputMediaDocument{
			Media:           filepath,
			Caption:         "Ваш файл! Наслаждайтесь чтением ❤️",
			MediaAttachment: bytes.NewBuffer(decompressed),
		},
	})
	if err != nil {
		b.SendDocument(ctx, &bot.SendDocumentParams{
			ChatID: update.Message.Chat.ID,
			ReplyParameters: &models.ReplyParameters{
				ChatID:    update.Message.Chat.ID,
				MessageID: update.Message.ID,
			},
			Document: &models.InputFileUpload{
				Data:     bytes.NewReader(decompressed),
				Filename: resp.Filename,
			},
		})
	}

	/*b.SendDocument(ctx, &bot.SendDocumentParams{
		ChatID: update.Message.Chat.ID,
		ReplyParameters: &models.ReplyParameters{
			ChatID:    update.Message.Chat.ID,
			MessageID: update.Message.ID,
		},
		Document: &models.InputFileUpload{
			Data:     bytes.NewReader(decompressed),
			Filename: resp.Filename,
		},
	})*/

}

func (h *EventProcessor) getPDFHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	h.logger.Debug("getPDFHandler(): "+update.Message.Text, zap.Int64("user", update.Message.From.ID))

	// Должен быть формат: /getpdf <description> (description - название)
	// Я не буду добавлять inline клавиатуру здесь, будет просто команда /getpdf с описанием и все

	// Делаем запрос к сервису, чтобы он вернул нам pdf в виде массива байтов
	// Если у сервиса нет контента, он вернет ошибку
	// Если все гуд, то вернется массив байтов, котоырй мы здесь преобразуем в pdf

	// ДОБАВИТЬ СЖАТИЕ С ДВУХ СТОРОН ДЛЯ ЭКОНОМИИ МЕСТА/ТРАФИКА
}

// Удаляем пдф по описанию (принцип схож с /del для ссылок, тоже будем возвращать список и давать юзеру выбрать, че удалить)
func (h *EventProcessor) deletePDFHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	h.logger.Debug("getPDFHandler(): "+update.Message.Text, zap.Int64("user", update.Message.From.ID))

}

// Удаляем просто по команде все PDF из памяти
func (h *EventProcessor) deleteAllPDFSHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	h.logger.Debug("getPDFHandler(): "+update.Message.Text, zap.Int64("user", update.Message.From.ID))

}

func (h *EventProcessor) inlineKbHandler(ctx context.Context, b *bot.Bot, update *models.Update, links []*gen.Link, tag string) {
	if len(links) == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Нет сохраненных ссылок 😞",
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
		Text:        "Choose the link",
		ReplyMarkup: inline,
	})
	if err != nil {
		h.logger.Error("err while send message (keyboard)", zap.Error(err))
	}
}

func (h *EventProcessor) helpHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	h.logger.Debug("helpHandler(): "+update.Message.Text, zap.Int64("user", update.Message.From.ID))

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      helpMsg_RU,
		ParseMode: models.ParseModeMarkdown,
	})
}
