package main

import (
	"database/sql"
	"log/slog"
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
	"github.com/modaniru/streamer-notifier-telegram/pkg/router"
)

// TODO wraps errors
// TODO uris to consts
func main() {
	config.LoadConfig()

	LoggerConfigure(os.Getenv("LEVEL"))
	slog.Info("logger was successfuly loaded")
	// TODO to .env
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:1111/postgres?sslmode=disable")
	if err != nil{
		slog.Error("postgres open error", slog.String("error", err.Error()))
		os.Exit(1)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_SECRET"))
	bot.Debug = true
	if err != nil{
		slog.Error("create bot api error", slog.String("error", err.Error()))
		os.Exit(1)
	}

	twitchClient := client.NewTwitchClient(&http.Client{}, os.Getenv("TWITCH_CLIENT_ID"), os.Getenv("TWITCH_CLIENT_SECRET"))
	storage := storage.NewStorage(db)
	service := service.NewService(storage, twitchClient)
	router := telegram.NewMyRouter(router.NewRouter(), bot, service)
	telegramBot := telegram.NewTelegramBot(bot, router.InitRouter())
	httpServer := server.NewServer(http.DefaultServeMux, telegramBot, service.StreamerService)

	status := make(chan int)
	slog.Info("start listening telegram bot...")
	go telegramBot.Listen(status)
	slog.Info("start listening http server...")
	go httpServer.Start("8080", status)

	if i := <-status; i == 1 {
		slog.Error("http server or telegram bot stop working")
		os.Exit(1)
	}
}

func LoggerConfigure(level string){
	var handler slog.Handler
	switch level{
	case "prod":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
			AddSource: true,
		})
	case "dev":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
			AddSource: true,
		})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
			AddSource: true,
		})
	}
	log := slog.New(handler)
	slog.SetDefault(log)
}