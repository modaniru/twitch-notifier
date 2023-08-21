package main

import (
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/modaniru/streamer-notifier-telegram/internal/client"
	"github.com/modaniru/streamer-notifier-telegram/internal/config"
	"github.com/modaniru/streamer-notifier-telegram/internal/server"
	"github.com/modaniru/streamer-notifier-telegram/internal/telegram"
)

// TODO wraps errors
// TODO uris to consts
func main() {
	config.LoadConfig()
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_SECRET"))
	if err != nil {
		log.Fatal(err.Error())
	}
	twitchClient := client.NewTwitchClient(&http.Client{}, os.Getenv("TWITCH_CLIENT_ID"), os.Getenv("TWITCH_CLIENT_SECRET"))
	telegramBot := telegram.NewTelegramBot(bot, twitchClient)
	httpServer := server.NewServer(http.DefaultServeMux, telegramBot)
	status := make(chan int)

	go telegramBot.Listen(status)
	go httpServer.Start("8080", status)

	if i := <-status; i == 1 {
		log.Fatal("telegram server error")
	}
}
