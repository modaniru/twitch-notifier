package storage

import (
	"database/sql"
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/modaniru/streamer-notifier-telegram/internal/storage/repo"
)

type User interface {
	SaveUserChatId(chatId int) (int, error)
	DeleteUser(chatId int) error
	GetUser(chatId int) (int, error)
}

type Streamers interface {
	SaveStreamer(streamerId string) (int, error)
	GetStreamer(streamerId string) (int, error)
	DeleteStreamer(streamerId string) error
}

type Follows interface {
	GetStreamersIdByChatId(chatId int) ([]string, error)
	SaveFollow(chatId int, streamerId int) error
	GetCountOfFollows(streamerId string) (int, error)
	GetAllStreamerFollowers(streamerId string) ([]int, error)
	Unfollow(chatId int, streamerId int) error
}

type Storage struct {
	User
	Streamers
	Follows
}

func NewStorage(db *sql.DB) *Storage {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err.Error())
	}
	return &Storage{
		User:      repo.NewUserStorage(db),
		Streamers: repo.NewStreamerStorage(db),
		Follows:   repo.NewFollowStorage(db),
	}
}
