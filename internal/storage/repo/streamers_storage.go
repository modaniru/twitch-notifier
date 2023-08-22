package repo

import "database/sql"

type StreamerStorage struct{
	db *sql.DB
}

func NewStreamerStorage(db *sql.DB) *StreamerStorage{
	return &StreamerStorage{db: db}
}

func (s *StreamerStorage) SaveStreamer(streamerId string) (int, error){
	sql := "insert into streamers (streamer_id) values ($1) returning id;"
	stmt, err := s.db.Prepare(sql)
	if err != nil{
		return 0, err
	}
	var id int
	err = stmt.QueryRow(streamerId).Scan(&id)
	if err != nil{
		return 0, err
	}
	return id, nil
}

func (u *StreamerStorage) DeleteStreamer(streamerId string) error{
	sql := "delete from streamers where streamer_id = $1;"
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

func (u *StreamerStorage) GetStreamer(streamerId string) (int, error){
	sql := "select id from streamers where streamer_id = $1;"
	stmt, err := u.db.Prepare(sql)
	if err != nil{
		return 0, err
	}
	var id int
	err = stmt.QueryRow(streamerId).Scan(&id)
	if err != nil{
		return 0, err
	}
	return id, nil
}