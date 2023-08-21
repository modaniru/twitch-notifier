package services

import (
	"database/sql"
	"errors"

	"github.com/modaniru/streamer-notifier-telegram/internal/storage"
)

type UserService struct{
	userStorage storage.User
}

func NewUserService(userStorage storage.User) *UserService{
	return &UserService{userStorage: userStorage}
}

func (u *UserService) CreateNewUser(chatId int) error{
	_, err := u.userStorage.GetUser(chatId)
	if err == nil{
		return nil
	}
	if errors.Is(err, sql.ErrNoRows){
		_, err := u.userStorage.SaveUserChatId(chatId)
		return err
	}
	return err
}