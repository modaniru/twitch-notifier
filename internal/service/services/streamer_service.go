package services

import (
	"database/sql"
	"errors"
	"os"
	"strconv"

	"github.com/modaniru/streamer-notifier-telegram/internal/client"
	"github.com/modaniru/streamer-notifier-telegram/internal/entity"
	"github.com/modaniru/streamer-notifier-telegram/internal/storage"
)

type StreamerService struct {
	followStorage   storage.Follows
	streamerStorage storage.Streamers
	twitchClient    *client.TwitchClient
}

func NewStreamerService(streamerStorage storage.Streamers, followStorage storage.Follows, twitchClient *client.TwitchClient) *StreamerService {
	return &StreamerService{
		streamerStorage: streamerStorage,
		followStorage:   followStorage,
		twitchClient:    twitchClient,
	}
}

// Return client.ErrStreamerNotFound if streamer was not found
func (s *StreamerService) SaveFollow(login string, chatId int) error {
	id, err := s.twitchClient.GetUserIdByLogin(login)
	if err != nil {
		return err
	}
	streamerId, err := s.streamerStorage.GetStreamer(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			streamerId, err = s.streamerStorage.SaveStreamer(id)
			if err != nil {
				return err
			}
			err = s.twitchClient.RegisterStreamWebhook("https://"+os.Getenv("DOMAIN")+"/stream-online", id)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	err = s.followStorage.SaveFollow(chatId, streamerId)
	if err != nil {
		return err
	}
	return nil
}

// Must return slice of nicknames
func (s *StreamerService) GetUserFollows(chatId int) ([]entity.UserInfo, error) {
	res, err := s.followStorage.GetStreamersIdByChatId(chatId)
	if err != nil {
		return nil, err
	}
	users, err := s.twitchClient.GetUsersByUserId(res)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *StreamerService) GetStreamerUsers(streamerId string) ([]int, error) {
	res, err := s.followStorage.GetAllStreamerFollowers(streamerId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *StreamerService) Unfollow(chatId int, login string) error {
	streamerIdString, err := s.twitchClient.GetUserIdByLogin(login)
	if err != nil {
		return err
	}
	id, err := strconv.Atoi(streamerIdString)
	if err != nil {
		return err
	}
	return s.followStorage.Unfollow(chatId, id)
}
