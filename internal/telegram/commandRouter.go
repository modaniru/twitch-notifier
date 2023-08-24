package telegram

import (
	"fmt"
	log "log/slog"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/modaniru/streamer-notifier-telegram/internal/service"
	"github.com/modaniru/streamer-notifier-telegram/internal/telegram/router"
)

type MyRouter struct{
	router *router.CommandRouter
	bot *tgbotapi.BotAPI
	service *service.Service
}

const(
	errorMsg = "Ошибка! 😶"
)

func NewMyRouter(router *router.CommandRouter, bot *tgbotapi.BotAPI, service *service.Service) *MyRouter{
	return &MyRouter{
		router: router,
		bot: bot,
		service: service,
	}
}

func SendMessage(bot *tgbotapi.BotAPI, message string, chatId int64) {
	msg := tgbotapi.NewMessage(chatId, message)
	_, err := bot.Send(msg)
	log.Info(fmt.Sprintf("send message to %d", chatId))
	if err != nil{
		log.Error("send message error", log.String("error", err.Error()))
	}
}

func (m *MyRouter) InitRouter() *router.CommandRouter{
	m.router.AddCommand("/start", router.Command{
		ArgumentsCount: 0,
		CommandHandler: m.StartCommand,
	})
	m.router.AddCommand("/ping", router.Command{
		ArgumentsCount: 0,
		CommandHandler: m.PingCommand,
	})
	m.router.AddCommand("/get", router.Command{
		ArgumentsCount: 0,
		CommandHandler: m.GetStreamers,
	})
	m.router.AddCommand("/add", router.Command{
		ArgumentsCount: 1,
		CommandHandler: m.AddStreamer,
	})

	return m.router
}

func (m *MyRouter) StartCommand(message *tgbotapi.Message) {
	chatId := message.From.ID
	textMessage := "Привет! 🥳\nДавай отслеживать твоих любимых стримеров! 😉"
	err := m.service.UserService.CreateNewUser(int(chatId))
	if err != nil{
		SendMessage(m.bot, errorMsg, chatId)
		return
	}
	SendMessage(m.bot, textMessage, chatId)
}

func (m *MyRouter) PingCommand(message *tgbotapi.Message){
	chatId := message.From.ID
	textMessage := "pong"
	SendMessage(m.bot, textMessage, chatId)
}

func (m *MyRouter) AddStreamer(message *tgbotapi.Message){
	chatId := message.From.ID
	streamerLogin := strings.Split(message.Text, " ")[1]
	userId, err := m.service.UserService.GetUser(int(chatId))
	if err != nil{
		log.Error("get user error", log.String("error", err.Error()))
		SendMessage(m.bot, errorMsg, chatId)
		return
	}
	err = m.service.StreamerService.SaveFollow(strings.ToLower(streamerLogin), userId)
	if err != nil{
		log.Error("save follow error", log.String("error", err.Error()))
		SendMessage(m.bot, errorMsg, chatId)
	} else {
		SendMessage(m.bot, "Успешно! 🥳", chatId)
	}
}

func (m *MyRouter) GetStreamers(message *tgbotapi.Message){
	chatId := message.From.ID
	res, err := m.service.StreamerService.GetUserFollows(int(chatId))
	if len(res) == 0{
		SendMessage(m.bot, "Вы не отслеживаете ни одного стримера! 😔", chatId)
		return
	}
	msg := ""
	for i, u := range res{
		msg += fmt.Sprintf("%d. %s\n", i + 1, u.DisplayName)
	}
	if err != nil{
		log.Error("get streamers error", log.String("error", err.Error()))
		SendMessage(m.bot, errorMsg, chatId)
	} else {
		SendMessage(m.bot, msg, chatId)
	}
}