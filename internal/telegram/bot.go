package telegram

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/modaniru/streamer-notifier-telegram/internal/entity"
	"github.com/modaniru/streamer-notifier-telegram/internal/service"
)

type TelegramBot struct {
	bot *tgbotapi.BotAPI
	service *service.Service
}

func NewTelegramBot(bot *tgbotapi.BotAPI, service *service.Service) *TelegramBot {
	return &TelegramBot{bot: bot, service: service}
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
				id := int(update.Message.From.ID)
				nickname := strings.Split(text, " ")[1]
				chatId, err := t.service.UserService.GetUser(id)
				if err != nil{
					t.SendMessage(err.Error(), int64(id))
					continue
				}
				err = t.service.StreamerService.SaveFollow(strings.ToLower(nickname), chatId)
				if err != nil{
					t.SendMessage(err.Error(), int64(id))
				} else {
					t.SendMessage("Успешно", int64(id))
				}
			} else if strings.HasPrefix(text, "/get "){
				id := int(update.Message.From.ID)
				res, err := t.service.StreamerService.GetUserFollows(id)
				if err != nil{
					t.SendMessage(err.Error(), int64(id))
				} else {
					t.SendMessage(fmt.Sprintf("%v", res), int64(id))
				}
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
