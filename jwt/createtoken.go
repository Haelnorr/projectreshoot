package jwt

import (
	"time"

	"projectreshoot/db"
	"projectreshoot/server"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

// Generates an access token for the provided user, using the variables set
// in the config object
func GenerateAccessToken(
	config *server.Config,
	user *db.User,
	fresh bool,
) (string, error) {
	issuedAt := time.Now().Unix()
	expiresAt := issuedAt + (config.AccessTokenExpiry * 60)
	var freshExpiresAt int64
	if fresh {
		freshExpiresAt = issuedAt + (config.TokenFreshTime * 60)
	} else {
		freshExpiresAt = issuedAt
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":   config.TrustedHost,
			"sub":   user.ID,
			"aud":   config.TrustedHost,
			"iat":   issuedAt,
			"exp":   expiresAt,
			"fresh": freshExpiresAt,
		})

	signedToken, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", errors.Wrap(err, "token.SignedString")
	}
	return signedToken, nil
}
