package storage

import (
	"database/sql"
	"log"

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
}

type Storage struct {
	User
	Streamers
	Follows
}

func NewStorage(db *sql.DB) *Storage {
	sql := `
	create table if not exists users(
		id serial primary key,
		chat_id int
	);
	
	create table if not exists streamers(
		id serial primary key,
		streamer_id varchar UNIQUE not null
	);
	
	create table if not exists follows(
		chat_id int REFERENCES users (id) on delete CASCADE,
		streamer_id int REFERENCES streamers (id) on delete CASCADE
	);`

	_, err := db.Exec(sql)
	if err != nil {
		log.Fatal(err.Error())
	}
	return &Storage{
		User:      repo.NewUserStorage(db),
		Streamers: repo.NewStreamerStorage(db),
		Follows:   repo.NewFollowStorage(db),
	}
}
