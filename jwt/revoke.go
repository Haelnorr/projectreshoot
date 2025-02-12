package jwt

import (
	"database/sql"

	"github.com/pkg/errors"
)

// Revoke a token by adding it to the database
func RevokeToken(conn *sql.DB, t Token) error {
	jti := t.GetJTI()
	exp := t.GetEXP()
	query := `INSERT INTO jwtblacklist (jti, exp) VALUES (?, ?)`
	_, err := conn.Exec(query, jti, exp)
	if err != nil {
		return errors.Wrap(err, "conn.Exec")
	}
	return nil
}

// Check if a token has been revoked. Returns true if not revoked.
func CheckTokenNotRevoked(conn *sql.DB, t Token) (bool, error) {
	jti := t.GetJTI()
	query := `SELECT 1 FROM jwtblacklist WHERE jti = ? LIMIT 1`
	rows, err := conn.Query(query, jti)
	defer rows.Close()
	if err != nil {
		return false, errors.Wrap(err, "conn.Exec")
	}
	revoked := rows.Next()
	return !revoked, nil
}
