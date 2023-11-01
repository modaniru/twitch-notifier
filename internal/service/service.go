package service

import (
	"github.com/modaniru/streamer-notifier-telegram/internal/client"
	"github.com/modaniru/streamer-notifier-telegram/internal/entity"
	"github.com/modaniru/streamer-notifier-telegram/internal/service/services"
	"github.com/modaniru/streamer-notifier-telegram/internal/storage"
)

type UserService interface {
	CreateNewUser(chatId int) error
	GetUser(chatId int) (int, error)
}

type StreamerService interface {
	GetUserFollows(chatId int) ([]entity.UserInfo, error)
	SaveFollow(login string, chatId int) error
	GetStreamerUsers(streamerId string) ([]int, error)
}

type Service struct {
	UserService
	StreamerService
}

func NewService(storage *storage.Storage, twitchClient *client.TwitchClient) *Service {
	return &Service{
		UserService:     services.NewUserService(storage.User),
		StreamerService: services.NewStreamerService(storage.Streamers, storage.Follows, twitchClient),
	}
}
