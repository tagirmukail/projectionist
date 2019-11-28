package db

const (
	CREATE_TBL_USERS = `
create table if not exists users
(
	id INTEGER
		constraint users_pk
			primary key autoincrement,
	username TEXT(255) not null,
	password TEXT(500) not null,
	role     int default 0,
	deleted  int default 0
);

create unique index if not exists users_username_uindex
    on users (username);
`
	CREATE_TBL_SERVICE = `
create table if not exists services
(
	id INTEGER
		constraint services_pk
			primary key autoincrement,
	name    TEXT(255) not null,
	link    TEXT(255) not null,
	Token   TEXT(255) not null,
	frequency int default 100,
	status  int default 1,
	deleted int default 0
);

create unique index if not exists services_name_uindex
    on services (name);
`
)
