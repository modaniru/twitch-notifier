package service

import (
	"github.com/modaniru/streamer-notifier-telegram/internal/service/services"
	"github.com/modaniru/streamer-notifier-telegram/internal/storage"
)

type UserService interface{
	CreateNewUser(chatId int) error
}

type Service struct{
	UserService
}

func NewService(storage *storage.Storage) *Service{
	return &Service{
		UserService: services.NewUserService(storage.User),
	}
}