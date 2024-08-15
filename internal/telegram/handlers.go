package telegram

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	log.Println("new update: ", update.Message.Text)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      "Hello",
		ParseMode: models.ParseModeMarkdown,
	})
}

func saveLinkHandler(ctx context.Context, b *bot.Bot, update *models.Update) {

}

func deleteLinkHandler(ctx context.Context, b *bot.Bot, update *models.Update) {

}

func getLinksHandler(ctx context.Context, b *bot.Bot, update *models.Update) {

}
