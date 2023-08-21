package storage

import (
	"database/sql"
	"log"

	"github.com/modaniru/streamer-notifier-telegram/internal/storage/repo"
)

type User interface{
	SaveUserChatId(chatId int) (int, error)
	DeleteUser(chatId int) error
	GetUser(chatId int) (int, error)
}

type Streamers interface{
	SaveStreamer(streamerId string) (int, error)
	DeleteStreamer(streamerId string) error
}

type Follows interface{
	SaveFollow(chatId int, streamerId string) error
	GetCountOfFollows(streamerId string) (int, error)
	GetAllStreamerFollowers(streamerId string) ([]int, error)
}

type Storage struct{
	User
	Streamers
	Follows
}

func NewStorage(db *sql.DB) *Storage{
	sql := `create table if exists users(
		id serial primary key,
		chat_id int
	);
	
	create table if exists streamers(
		id serial primary key,
		streamer_id varchar UNIQUE not null
	);
	
	create table if exists follow(
		user_id int,
		streamer_id int,
		foreign key (user_id) REFERENCES users (id) on delete CASCADE,
		foreign key (streamer_id) REFERENCES streamers (id) on delete CASCADE
	);`

	_, err := db.Exec(sql)
	if err != nil{
		log.Fatal(err.Error())
	}
	return &Storage{
		User: repo.NewUserStorage(db),
		Streamers: repo.NewStreamerStorage(db),
		Follows: repo.NewFollowStorage(db),
	}
}