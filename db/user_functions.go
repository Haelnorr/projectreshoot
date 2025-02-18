package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

// Creates a new user in the database and returns a pointer
func CreateNewUser(
	ctx context.Context,
	tx *SafeTX,
	username string,
	password string,
) (*User, error) {
	query := `INSERT INTO users (username) VALUES (?)`
	_, err := tx.Exec(ctx, query, username)
	if err != nil {
		return nil, errors.Wrap(err, "tx.Exec")
	}
	user, err := GetUserFromUsername(ctx, tx, username)
	if err != nil {
		return nil, errors.Wrap(err, "GetUserFromUsername")
	}
	err = user.SetPassword(ctx, tx, password)
	if err != nil {
		return nil, errors.Wrap(err, "user.SetPassword")
	}
	return user, nil
}

// Fetches data from the users table using "WHERE column = 'value'"
func fetchUserData(
	ctx context.Context,
	tx *SafeTX,
	column string,
	value interface{},
) (*sql.Rows, error) {
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
	rows, err := tx.Query(ctx, query, value)
	if err != nil {
		return nil, errors.Wrap(err, "tx.Query")
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
func GetUserFromUsername(ctx context.Context, tx *SafeTX, username string) (*User, error) {
	rows, err := fetchUserData(ctx, tx, "username", username)
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
func GetUserFromID(ctx context.Context, tx *SafeTX, id int) (*User, error) {
	rows, err := fetchUserData(ctx, tx, "id", id)
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
func CheckUsernameUnique(ctx context.Context, tx *SafeTX, username string) (bool, error) {
	query := `SELECT 1 FROM users WHERE username = ? COLLATE NOCASE LIMIT 1`
	rows, err := tx.Query(ctx, query, username)
	if err != nil {
		return false, errors.Wrap(err, "tx.Query")
	}
	defer rows.Close()
	taken := rows.Next()
	return !taken, nil
}
