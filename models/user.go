package models

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type Role int

const (
	Admin Role = iota
	SuperAdmin
)

type User struct {
	ID       int    `json:"id";db:"id"`
	Username string `json:"username";db:"username"`
	Role     int    `json:"role";db:"role"`
	Password string `db:"password"`
	Token    string `json:"token"`
	Deleted  int    `json:"deleted";db:"deleted"`
}

func (u *User) Validate() error {
	if u.Username == "" && len(u.Username) > 255 {
		return fmt.Errorf("Invalid username")
	}
	if u.Password == "" && len(u.Password) > 500 {
		return fmt.Errorf("Invalid password")
	}

	switch u.Role {
	case int(Admin):
		break
	case int(SuperAdmin):
		break
	default:
		return fmt.Errorf("Invalid role")
	}

	return nil
}

func (u *User) IsExist(db *sql.DB) (error, bool) {
	var username string
	var err = db.QueryRow("SELECT username FROM users where username=?", u.Username).Scan(&username)
	if err != nil && err != sql.ErrNoRows {
		return err, false
	}

	if username == "" {
		return fmt.Errorf("User with username %s not exist", u.Username), false
	}

	return nil, true
}

func (u *User) Save(db *sql.DB) error {
	passwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)
	if err != nil {
		return err
	}

	result, err := db.Exec(
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

func (u *User) Count(db *sql.DB) (int, error) {
	var count int
	var err = db.QueryRow("SELECT count(id) FROM users").Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (u *User) GetByName(db *sql.DB, username string) error {
	return db.QueryRow("SELECT id, username, password, role, deleted FROM users WHERE username=?", username).Scan(
		&u.ID,
		&u.Username,
		&u.Password,
		&u.Role,
		&u.Deleted,
	)
}

func (u *User) GetByID(db *sql.DB, id int64) error {
	return db.QueryRow("SELECT id, username, password, role, deleted FROM users WHERE id=?", id).Scan(
		&u.ID,
		&u.Username,
		&u.Password,
		&u.Role,
		&u.Deleted,
	)
}

func (u *User) Pagination(db *sql.DB, start, end int) ([]Model, error) {
	var result []Model

	raws, err := db.Query("SELECT id, username, role, deleted FROM users ORDER BY id ASC limit ?, ?", start, end)
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

func (u *User) Update(db *sql.DB) error {
	query := "UPDATE users SET "

	if u.Username != "" {
		query += fmt.Sprintf("username='%s', ", u.Username)
	}

	query += fmt.Sprintf("role=%d ", u.Role)

	query += fmt.Sprintf("WHERE id=%d", u.ID)

	_, err := db.Exec(
		query,
	)

	return err
}

func (u *User) Delete(db *sql.DB, id int) error {
	query := fmt.Sprintf("DELETE FROM users WHERE id=%d", id)

	_, err := db.Exec(query)
	return err
}
