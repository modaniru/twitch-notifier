package repo

import "database/sql"

type StreamerStorage struct{
	db *sql.DB
}

func NewStreamerStorage(db *sql.DB) *StreamerStorage{
	return &StreamerStorage{db: db}
}

func (s *StreamerStorage) SaveStreamer(streamerId string) (int, error){
	sql := "insert into streamers (streamer_id) values (?);"
	stmt, err := s.db.Prepare(sql)
	if err != nil{
		return 0, err
	}
	response, err := stmt.Exec(streamerId)
	if err != nil{
		return 0, err
	}
	id, err := response.LastInsertId()
	if err != nil{
		return 0, err
	}
	return int(id), err
}

func (u *StreamerStorage) DeleteStreamer(streamerId string) error{
	sql := "delete from streamers where streamer_id = ?;"
	stmt, err := u.db.Prepare(sql)
	if err != nil{
		return err
	}
	_, err = stmt.Exec(streamerId)
	if err != nil{
		return err
	}
	return err
}