package repo

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	mye "github.com/modaniru/streamer-notifier-telegram/internal/errors"
)

type FollowStorage struct {
	db *sql.DB
}

func NewFollowStorage(db *sql.DB) *FollowStorage {
	return &FollowStorage{db: db}
}

func (f *FollowStorage) SaveFollow(chatId int, streamerId int) error {
	op := "SaveFollow"
	sql := "insert into follows (chat_id, streamer_id) values ($1, $2);"
	stmt, err := f.db.Prepare(sql)
	if err != nil {
		return fmt.Errorf("prepare %s error: %w", op, err)
	}
	_, err = stmt.Exec(chatId, streamerId)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return mye.ErrFollowAlreadyExists
		}
		return fmt.Errorf("exec %s error: %w", op, err)
	}
	return nil
}

func (f *FollowStorage) Unfollow(chatId int, streamerId int) error {
	op := "Unfollow"
	query := `DELETE FROM follows as f where chat_id = (select id from users where chat_id = $1) and streamer_id = (select id from streamers where streamer_id = $2);`
	stmt, err := f.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("prepare %s error: %w", op, err)
	}
	res, err := stmt.Exec(chatId, streamerId)
	// a := 0
	// err = row.Scan(&a, &a)
	if i, err := res.RowsAffected(); err != nil || i == 0 {
		if err != nil {
			return err
		}
		return mye.ErrFollowNotFound
	}
	if err != nil {
		return fmt.Errorf("exec %s error: %w", op, err)
	}
	return nil
}

func (f *FollowStorage) GetCountOfFollows(streamerId string) (int, error) {
	op := "GetCountOfFollows"
	sql := "select count() from streamers inner join follows on streamers.id = follows.streamer_id group by streamers.streamer_id;"
	stmt, err := f.db.Prepare(sql)
	if err != nil {
		return 0, fmt.Errorf("prepare %s error: %w", op, err)
	}
	var count int
	err = stmt.QueryRow().Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("exec %s error: %w", op, err)
	}
	return count, nil
}
func (f *FollowStorage) GetAllStreamerFollowers(streamerId string) ([]int, error) {
	op := "GetAllStreamerFollowers"
	sql := "select u.chat_id from follows as f inner join streamers as s on f.streamer_id = s.id inner join users as u on f.chat_id = u.id where s.streamer_id = $1;"
	stmt, err := f.db.Prepare(sql)
	if err != nil {
		return nil, fmt.Errorf("prepare %s error: %w", op, err)
	}
	rows, err := stmt.Query(streamerId)
	if err != nil {
		return nil, fmt.Errorf("exec %s error: %w", op, err)
	}
	rows.Columns()
	defer rows.Close()
	res := make([]int, 0)
	for rows.Next() {
		var chatId int
		err := rows.Scan(&chatId)
		if err != nil {
			return nil, fmt.Errorf("scan %s error: %w", op, err)
		}
		res = append(res, chatId)
	}
	return res, nil
}

func (f *FollowStorage) GetStreamersIdByChatId(chatId int) ([]string, error) {
	op := "GetStreamersIdByChatId"
	sql := "select s.streamer_id from users as u inner join follows as f on u.id = f.chat_id inner join streamers as s on f.streamer_id = s.id where u.chat_id = $1;"
	stmt, err := f.db.Prepare(sql)
	if err != nil {
		return nil, fmt.Errorf("prepare %s error: %w", op, err)
	}
	rows, err := stmt.Query(chatId)
	if err != nil {
		return nil, fmt.Errorf("exec %s error: %w", op, err)
	}
	defer rows.Close()
	res := make([]string, 0)
	for rows.Next() {
		var streamerId string
		err := rows.Scan(&streamerId)
		if err != nil {
			return nil, fmt.Errorf("scan %s error: %w", op, err)
		}
		res = append(res, streamerId)
	}
	return res, nil
}
