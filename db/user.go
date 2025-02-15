package db

import (
	"database/sql"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID            int    // Integer ID (index primary key)
	Username      string // Username (unique)
	Password_hash string // Bcrypt password hash
	Created_at    int64  // Epoch timestamp when the user was added to the database
	Bio           string // Short byline set by the user
}

// Uses bcrypt to set the users Password_hash from the given password
func (user *User) SetPassword(conn *sql.DB, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, "bcrypt.GenerateFromPassword")
	}
	user.Password_hash = string(hashedPassword)
	query := `UPDATE users SET password_hash = ? WHERE id = ?`
	_, err = conn.Exec(query, user.Password_hash, user.ID)
	if err != nil {
		return errors.Wrap(err, "conn.Exec")
	}
	return nil
}

// Uses bcrypt to check if the given password matches the users Password_hash
func (user *User) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password_hash), []byte(password))
	if err != nil {
		return errors.Wrap(err, "bcrypt.CompareHashAndPassword")
	}
	return nil
}

// Change the user's username
func (user *User) ChangeUsername(conn *sql.DB, newUsername string) error {
	query := `UPDATE users SET username = ? WHERE id = ?`
	_, err := conn.Exec(query, newUsername, user.ID)
	if err != nil {
		return errors.Wrap(err, "conn.Exec")
	}
	return nil
}
