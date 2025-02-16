package db

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

// Creates a new user in the database and returns a pointer
func CreateNewUser(conn *sql.DB, username string, password string) (*User, error) {
	query := `INSERT INTO users (username) VALUES (?)`
	_, err := conn.Exec(query, username)
	if err != nil {
		return nil, errors.Wrap(err, "conn.Exec")
	}
	user, err := GetUserFromUsername(conn, username)
	if err != nil {
		return nil, errors.Wrap(err, "GetUserFromUsername")
	}
	err = user.SetPassword(conn, password)
	if err != nil {
		return nil, errors.Wrap(err, "user.SetPassword")
	}
	return user, nil
}

// Fetches data from the users table using "WHERE column = 'value'"
func fetchUserData(conn *sql.DB, column string, value interface{}) (*sql.Rows, error) {
	query := fmt.Sprintf(
		`SELECT 
            id,
            username, 
            password_hash, 
            created_at,
            bio
        FROM users 
	    WHERE %s = ? COLLATE NOCASE LIMIT 1`,
		column,
	)
	rows, err := conn.Query(query, value)
	if err != nil {
		return nil, errors.Wrap(err, "conn.Query")
	}
	return rows, nil
}

// Scan the next row into the provided user pointer. Calls rows.Next() and
// assumes only row in the result. Providing a rows object with more than 1
// row may result in undefined behaviour.
func scanUserRow(user *User, rows *sql.Rows) error {
	for rows.Next() {
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Password_hash,
			&user.Created_at,
			&user.Bio,
		)
		if err != nil {
			return errors.Wrap(err, "rows.Scan")
		}
	}
	return nil
}

// Queries the database for a user matching the given username.
// Query is case insensitive
func GetUserFromUsername(conn *sql.DB, username string) (*User, error) {
	rows, err := fetchUserData(conn, "username", username)
	if err != nil {
		return nil, errors.Wrap(err, "fetchUserData")
	}
	defer rows.Close()
	var user User
	err = scanUserRow(&user, rows)
	if err != nil {
		return nil, errors.Wrap(err, "scanUserRow")
	}
	return &user, nil
}

// Queries the database for a user matching the given ID.
func GetUserFromID(conn *sql.DB, id int) (*User, error) {
	rows, err := fetchUserData(conn, "id", id)
	if err != nil {
		return nil, errors.Wrap(err, "fetchUserData")
	}
	defer rows.Close()
	var user User
	err = scanUserRow(&user, rows)
	if err != nil {
		return nil, errors.Wrap(err, "scanUserRow")
	}
	return &user, nil
}

// Checks if the given username is unique. Returns true if not taken
func CheckUsernameUnique(conn *sql.DB, username string) (bool, error) {
	query := `SELECT 1 FROM users WHERE username = ? COLLATE NOCASE LIMIT 1`
	rows, err := conn.Query(query, username)
	if err != nil {
		return false, errors.Wrap(err, "conn.Query")
	}
	defer rows.Close()
	taken := rows.Next()
	return !taken, nil
}
