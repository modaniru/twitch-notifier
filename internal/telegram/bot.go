package telegram

import (
	"errors"
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/modaniru/streamer-notifier-telegram/internal/client"
	"github.com/modaniru/streamer-notifier-telegram/internal/entity"
	"github.com/modaniru/streamer-notifier-telegram/internal/service"
)

type TelegramBot struct {
	bot *tgbotapi.BotAPI
	twitchClient *client.TwitchClient
	service *service.Service
}

func NewTelegramBot(bot *tgbotapi.BotAPI, twitchClient *client.TwitchClient, service *service.Service) *TelegramBot {
	return &TelegramBot{bot: bot, twitchClient: twitchClient, service: service}
}

func (t *TelegramBot) Listen(status chan int) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.bot.GetUpdatesChan(u)

	log.Println("start telegram bot")
	for update := range updates {
		if update.Message != nil { // If we got a message
			text := update.Message.Text
			if strings.HasPrefix(text, "/start"){
				id := int(update.Message.From.ID)
				err := t.service.CreateNewUser(id)
				if err != nil{
					t.SendMessage("Ошибка!", int64(id))
					continue
				}
				t.SendMessage("Привет! Начнем отслеживать твоих любимых стримеров?", int64(id))
			} else if strings.HasPrefix(text, "/add "){
				id, err := t.twitchClient.GetUserIdByLogin(strings.Split(text, " ")[1])
				if err != nil{
					if errors.Is(err, client.ErrStreamerNotFound){
						t.SendMessage("Стример не найден!", update.Message.From.ID)
						continue
					}
					t.SendMessage("Internal server error", update.Message.From.ID)
					continue
				}
				t.SendMessage(fmt.Sprintf("Стример найден его user_id %s", id), update.Message.From.ID)
				continue
			}
		}
	}
	status <- 1
}

func (t *TelegramBot) SendMessage(message string, chatId int64) {
	msg := tgbotapi.NewMessage(451819182, message)
	t.bot.Send(msg)
}

func (t *TelegramBot) SendNotification(notification entity.StreamOnlineNotification, chatId int64){
	login := notification.Event.BroadcasterUserLogin
	uri := fmt.Sprintf("twitch.tv/%s", login)
	message := fmt.Sprintf("🔴 %s запустил стрим\n%s", login, uri)
	t.SendMessage(message, chatId)
}
