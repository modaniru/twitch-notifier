create table users(
    id int primary key AUTOINCREMENT,
    chat_id int
);

create table streamers(
    id int primary key AUTOINCREMENT,
    streamer_id varchar UNIQUE not null
);

create table follow(
    user_id int,
    streamer_id int,
    foreign key (user_id) REFERENCES users (id) on delete CASCADE,
    foreign key (streamer_id) REFERENCES streamers (id) on delete CASCADE
);