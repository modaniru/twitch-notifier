package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
	"github.com/modaniru/streamer-notifier-telegram/internal/client"
	"github.com/modaniru/streamer-notifier-telegram/internal/config"
	"github.com/modaniru/streamer-notifier-telegram/internal/server"
	"github.com/modaniru/streamer-notifier-telegram/internal/service"
	"github.com/modaniru/streamer-notifier-telegram/internal/storage"
	"github.com/modaniru/streamer-notifier-telegram/internal/telegram"
)

// TODO wraps errors
// TODO uris to consts
func main() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:1111/postgres?sslmode=disable")
	if err != nil{
		log.Fatal(err.Error())
	}
	storage := storage.NewStorage(db)
	service := service.NewService(storage)
	config.LoadConfig()
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_SECRET"))
	if err != nil {
		log.Fatal(err.Error())
	}
	twitchClient := client.NewTwitchClient(&http.Client{}, os.Getenv("TWITCH_CLIENT_ID"), os.Getenv("TWITCH_CLIENT_SECRET"))
	telegramBot := telegram.NewTelegramBot(bot, twitchClient, service)
	httpServer := server.NewServer(http.DefaultServeMux, telegramBot)
	status := make(chan int)

	go telegramBot.Listen(status)
	go httpServer.Start("8080", status)

	if i := <-status; i == 1 {
		log.Fatal("telegram server error")
	}
}
