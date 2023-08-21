package repo

import "database/sql"

type FollowStorage struct{
	db *sql.DB
}

func NewFollowStorage(db *sql.DB) *FollowStorage{
	return &FollowStorage{db: db}
}

func (f *FollowStorage) SaveFollow(chatId int, streamerId string) error{
	sql := "insert into follows (chat_id, streamer_id) values (?, ?);"
	stmt, err := f.db.Prepare(sql)
	if err != nil{
		return err
	}
	_, err = stmt.Exec(chatId, streamerId)
	if err != nil{
		return err
	}
	return nil
}

func (f *FollowStorage) GetCountOfFollows(streamerId string) (int, error){
	sql := "select count() from streamers inner join follow on streamers.id = follow.streamer_id group by streamers.streamer_id;"
	stmt, err := f.db.Prepare(sql)
	if err != nil{
		return 0, err
	}
	var count int
	err = stmt.QueryRow().Scan(&count)
	if err != nil{
		return 0, err
	}
	return count, nil
}
func (f *FollowStorage) GetAllStreamerFollowers(streamerId string) ([]int, error){
	sql := "select u.chat_id from users as u inner join follow as f on u.id = f.user_id inner join streamers as s on f.streamer_id = s.streamer_id where f.streamer_id = ?;"
	stmt, err := f.db.Prepare(sql)
	if err != nil{
		return nil, err
	}
	rows, err := stmt.Query(streamerId)
	if err != nil{
		return nil, err
	}
	defer rows.Close()
	res := make([]int, 0)
	for rows.Next(){
		var chatId int
		err := rows.Scan(&chatId)
		if err != nil{
			return nil, err
		}
		res = append(res, chatId)
	}
	return res, nil
}