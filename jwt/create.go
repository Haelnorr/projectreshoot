package jwt

import (
	"time"

	"projectreshoot/db"
	"projectreshoot/server"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Generates an access token for the provided user
func GenerateAccessToken(
	config *server.Config,
	user *db.User,
	fresh bool,
	rememberMe bool,
) (string, error) {
	issuedAt := time.Now().Unix()
	expiresAt := issuedAt + (config.AccessTokenExpiry * 60)
	var freshExpiresAt int64
	if fresh {
		freshExpiresAt = issuedAt + (config.TokenFreshTime * 60)
	} else {
		freshExpiresAt = issuedAt
	}
	var ttl string
	if rememberMe {
		ttl = "exp"
	} else {
		ttl = "session"
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":   config.TrustedHost,
			"scope": "access",
			"ttl":   ttl,
			"iat":   issuedAt,
			"exp":   expiresAt,
			"fresh": freshExpiresAt,
			"sub":   user.ID,
		})

	signedToken, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", errors.Wrap(err, "token.SignedString")
	}
	return signedToken, nil
}

// Generates a refresh token for the provided user
func GenerateRefreshToken(
	config *server.Config,
	user *db.User,
	rememberMe bool,
) (string, error) {
	issuedAt := time.Now().Unix()
	expiresAt := issuedAt + (config.RefreshTokenExpiry * 60)
	var ttl string
	if rememberMe {
		ttl = "exp"
	} else {
		ttl = "session"
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":   config.TrustedHost,
			"scope": "refresh",
			"ttl":   ttl,
			"jti":   uuid.New(),
			"iat":   issuedAt,
			"exp":   expiresAt,
			"sub":   user.ID,
		})

	signedToken, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", errors.Wrap(err, "token.SignedString")
	}
	return signedToken, nil
}
