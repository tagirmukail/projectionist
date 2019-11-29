package models

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type Role int

const (
	_ Role = iota
	Admin
	SuperAdmin
)

type User struct {
	dbCtx    *sql.DB
	ID       int    `json:"id";db:"id"`
	Username string `json:"username";db:"username"`
	Role     Role   `json:"role";db:"role"`
	Password string `db:"password"`
	Token    string `json:"token"`
	Deleted  int    `json:"deleted";db:"deleted"`
}

func (u *User) SetDBCtx(iDB interface{}) error {
	db, ok := iDB.(*sql.DB)
	if !ok {
		return fmt.Errorf("%v is not sql.DB", iDB)
	}

	u.dbCtx = db

	return nil
}

func (u *User) Validate() error {
	if u.Username == "" || len(u.Username) > 255 {
		return fmt.Errorf("invalid username")
	}
	if u.Password == "" || len(u.Password) > 500 {
		return fmt.Errorf("invalid password")
	}

	switch u.Role {
	case Admin:
		break
	case SuperAdmin:
		break
	default:
		return fmt.Errorf("invalid role")
	}

	return nil
}

func (u *User) IsExistByName() (error, bool) {
	var username string
	var err = u.dbCtx.QueryRow("SELECT username FROM users where username=?", u.Username).Scan(&username)
	if err != nil && err != sql.ErrNoRows {
		return err, false
	}

	if username == "" {
		return nil, false
	}

	return nil, true
}

func (u *User) Save() error {
	passwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)
	if err != nil {
		return err
	}

	result, err := u.dbCtx.Exec(
		"INSERT INTO users (username, password, role) VALUES (?,?,?)",
		u.Username,
		string(passwd),
		u.Role,
	)
	if err != nil {
		return err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = int(lastInsertID)

	return nil
}

func (u *User) Count() (int, error) {
	var count int
	var err = u.dbCtx.QueryRow("SELECT count(id) FROM users").Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (u *User) GetByName(username string) error {
	return u.dbCtx.QueryRow("SELECT id, username, password, role, deleted FROM users WHERE username=?", username).Scan(
		&u.ID,
		&u.Username,
		&u.Password,
		&u.Role,
		&u.Deleted,
	)
}

func (u *User) GetByID(id int64) error {
	return u.dbCtx.QueryRow("SELECT id, username, password, role, deleted FROM users WHERE id=?", id).Scan(
		&u.ID,
		&u.Username,
		&u.Password,
		&u.Role,
		&u.Deleted,
	)
}

func (u *User) Pagination(start, end int) ([]Model, error) {
	var result []Model

	raws, err := u.dbCtx.Query(
		"SELECT id, username, role, deleted FROM users ORDER BY id ASC limit ?, ?", start, end)
	if err != nil {
		return result, err
	}

	for raws.Next() {
		var user = &User{}

		err = raws.Scan(&user.ID, &user.Username, &user.Role, &user.Deleted)
		if err != nil {
			return result, err
		}

		result = append(result, user)
	}

	return result, nil
}

func (u *User) Update(id int) error {
	var args []interface{}

	query := `UPDATE users SET `

	if u.Username != "" {
		query += `username=?, `
		args = append(args, u.Username)
	}

	if u.Role != Role(0) {
		query += `role=? `
		args = append(args, u.Role)
	}

	query += `WHERE id=?`
	args = append(args, id)

	res, err := u.dbCtx.Exec(
		query, args...,
	)
	if err != nil {
		return err
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsCount == 0 {
		return fmt.Errorf("user with id %d not updated", id)
	}

	return nil
}

func (u *User) Delete(id int) error {
	res, err := u.dbCtx.Exec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		return err
	}

	rowCount, err := res.LastInsertId()
	if err != nil {
		return err
	}

	if rowCount == 0 {
		return fmt.Errorf("user with id %v not deleted", id)
	}

	return err
}

func (u *User) GetID() int {
	return u.ID
}

func (u *User) SetID(id int) {
	u.ID = id
}

func (u *User) GetName() string {
	return u.Username
}

func (u *User) SetDeleted() {
	u.Deleted = 1
}

func (u *User) IsDeleted() bool {
	return u.Deleted > 0
}
