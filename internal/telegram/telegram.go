package telegram

import (
	api "github.com/0x0FACED/link-saver-telegram-bot/internal/grpc"
	"github.com/0x0FACED/link-saver-telegram-bot/internal/logger"
	"github.com/go-telegram/bot"
)

type TelegramBot struct {
	bot       *bot.Bot
	apiClient *api.APIClient
	logger    *logger.Logger
}

// func New(token string, apiHost string, conn *grpc.ClientConn) *TelegramBot {
// 	bot, err := bot.New(token)
// 	if err != nil {
// 		return nil
// 	}

// 	client := api.New(conn)

// }
func (b *TelegramBot) Start() {

}
