package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/modaniru/streamer-notifier-telegram/internal/entity"
	"github.com/modaniru/streamer-notifier-telegram/pkg/router"
)

type TelegramBot struct {
	bot    *tgbotapi.BotAPI
	router *router.CommandRouter
}

func NewTelegramBot(bot *tgbotapi.BotAPI, router *router.CommandRouter) *TelegramBot {
	return &TelegramBot{bot: bot, router: router}
}

// TODO add command sender abstraction
func (t *TelegramBot) Listen(status chan int) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := t.bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		err := t.router.Route(update.Message)
		if err != nil {
			SendMessage(t.bot, err.Error(), update.Message.From.ID)
		}
	}

	status <- 1
}

func (t *TelegramBot) SendNotification(notification entity.StreamOnlineNotification, chatId int64) {
	login := notification.Event.BroadcasterUserLogin
	uri := fmt.Sprintf("twitch.tv/%s", login)
	message := fmt.Sprintf("🔴 %s запустил стрим\n%s", login, uri)
	SendMessage(t.bot, message, chatId)
}
