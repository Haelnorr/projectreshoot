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
}

// Uses bcrypt to set the users Password_hash from the given password
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
	ra, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "result.RowsAffected")
	}
	if ra != 1 {
		return errors.New("Password was not updated")
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

// Queries the database for a user matching the given username.
// Query is case insensitive
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

// Queries the database for a user matching the given ID.
func GetUserFromID(conn *sql.DB, id int) (User, error) {
	query := `SELECT id, username, password_hash, created_at FROM users 
	          WHERE id = ?`
	rows, err := conn.Query(query, id)
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
