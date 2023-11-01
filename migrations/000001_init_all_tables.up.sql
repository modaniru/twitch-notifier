create table if not exists users(
	id serial primary key,
	chat_id int unique
);
	
create table if not exists streamers(
	id serial primary key,
	streamer_id varchar UNIQUE not null
);
	
create table if not exists follows(
	chat_id int REFERENCES users (id) on delete CASCADE,
	streamer_id int REFERENCES streamers (id) on delete CASCADE,
	UNIQUE(chat_id, streamer_id)
);