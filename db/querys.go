package db

const (
	CREATE_TBL_USERS = `
create table if not exists users
(
	id INTEGER
		constraint users_pk
			primary key autoincrement,
	username TEXT(255) not null,
	role     int default 0,
	deleted  int default 0
);

create unique index users_username_uindex
    on users (username);
`
	CREATE_TBL_SERVICE = `
create table if not exists services
(
	id INTEGER
		constraint services_pk
			primary key autoincrement,
	name    TEXT(255) not null,
	pid     int default 0,
	status  int default 0,
	deleted int default 0
);

create unique index services_name_uindex
    on services (name);
`
)