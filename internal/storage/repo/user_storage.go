package repo

import "database/sql"

type UserStorage struct{
	db *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage{
	return &UserStorage{db: db}
}

func (u *UserStorage) SaveUserChatId(chatId int) (int, error){
	sql := "insert into users (chat_id) values ($1) returning id;"
	stmt, err := u.db.Prepare(sql)
	if err != nil{
		return 0, err
	}
	var id int
	err = stmt.QueryRow(chatId).Scan(&id)
	if err != nil{
		return 0, err
	}
	return int(id), err
}

func (u *UserStorage) DeleteUser(chatId int) error{
	sql := "delete from users where chat_id = $1;"
	stmt, err := u.db.Prepare(sql)
	if err != nil{
		return err
	}
	_, err = stmt.Exec(chatId)
	if err != nil{
		return err
	}
	return err
}

func (u *UserStorage) GetUser(chatId int) (int, error){
	sql := "select id from users where chat_id = $1;"
	stmt, err := u.db.Prepare(sql)
	if err != nil{
		return 0, err
	}
	var id int
	err = stmt.QueryRow(chatId).Scan(&id)
	if err != nil{
		return 0, err
	}
	return int(id), err
}