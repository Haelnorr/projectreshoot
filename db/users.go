package db

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID            int
	Username      string
	Password_hash string
	Created_at    int64
}

func (user *User) SetPassword(conn *sql.DB, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, "bcrypt.GenerateFromPassword")
	}
	user.Password_hash = string(hashedPassword)
	query := `UPDATE users SET password_hash = ? WHERE id = ?`
	result, err := conn.Exec(query, user.Password_hash, user.ID)
	if err != nil {
		return errors.Wrap(err, "conn.Exec")
	}
	fmt.Println(result)
	return nil
}

func (user *User) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password_hash), []byte(password))
	if err != nil {
		return errors.Wrap(err, "bcrypt.CompareHashAndPassword")
	}
	return nil
}

func GetUserFromUsername(conn *sql.DB, username string) (User, error) {
	query := `SELECT id, username, password_hash, created_at FROM users 
	          WHERE username = ? COLLATE NOCASE`
	rows, err := conn.Query(query, username)
	if err != nil {
		return User{}, errors.Wrap(err, "conn.Query")
	}
	defer rows.Close()
	var user User
	for rows.Next() {
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Password_hash,
			&user.Created_at,
		)
		if err != nil {
			return User{}, errors.Wrap(err, "rows.Scan")
		}
	}
	return user, nil
}
