package repo

import (
	"database/sql"
	"fmt"
)

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{db: db}
}

func (u *UserStorage) SaveUserChatId(chatId int) (int, error) {
	op := "SaveUserChatId"
	sql := "insert into users (chat_id) values ($1) returning id;"
	stmt, err := u.db.Prepare(sql)
	if err != nil {
		return 0, fmt.Errorf("prepare %s error: %w", op, err)
	}
	var id int
	err = stmt.QueryRow(chatId).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("exec %s error: %w", op, err)
	}
	return int(id), err
}

func (u *UserStorage) DeleteUser(chatId int) error {
	op := "DeleteUser"

	sql := "delete from users where chat_id = $1;"
	stmt, err := u.db.Prepare(sql)
	if err != nil {
		return fmt.Errorf("prepare %s error: %w", op, err)
	}
	_, err = stmt.Exec(chatId)
	if err != nil {
		return fmt.Errorf("exec %s error: %w", op, err)
	}
	return err
}

func (u *UserStorage) GetUser(chatId int) (int, error) {
	op := "GetUser"
	sql := "select id from users where chat_id = $1;"
	stmt, err := u.db.Prepare(sql)
	if err != nil {
		return 0, fmt.Errorf("prepare %s error: %w", op, err)
	}
	var id int
	err = stmt.QueryRow(chatId).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("exec %s error: %w", op, err)
	}
	return int(id), err
}
