package jwt

import (
	"fmt"
	"projectreshoot/server"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Parse an access token and return a struct with all the claims. Does validation on
// all the claims, including checking if it is expired, has a valid issuer, and
// has the correct scope.
func ParseAccessToken(config *server.Config, tokenString string) (AccessToken, error) {
	claims, err := parseToken(config.SecretKey, tokenString)
	if err != nil {
		return AccessToken{}, errors.Wrap(err, "parseToken")
	}
	expiry, err := checkTokenExpired(claims["exp"])
	if err != nil {
		return AccessToken{}, errors.Wrap(err, "checkTokenExpired")
	}
	issuer, err := checkTokenIssuer(config.TrustedHost, claims["iss"])
	if err != nil {
		return AccessToken{}, errors.Wrap(err, "checkTokenIssuer")
	}
	ttl, err := getTokenTTL(claims["ttl"])
	if err != nil {
		return AccessToken{}, errors.Wrap(err, "getTokenTTL")
	}
	scope, err := getTokenScope(claims["scope"])
	if err != nil {
		return AccessToken{}, errors.Wrap(err, "getTokenScope")
	}
	if scope != "access" {
		return AccessToken{}, errors.New("Token is not an Access token")
	}
	issuedAt, err := getIssuedTime(claims["iat"])
	if err != nil {
		return AccessToken{}, errors.Wrap(err, "getIssuedTime")
	}
	subject, err := getTokenSubject(claims["sub"])
	if err != nil {
		return AccessToken{}, errors.Wrap(err, "getTokenSubject")
	}
	fresh, err := getFreshTime(claims["fresh"])
	if err != nil {
		return AccessToken{}, errors.Wrap(err, "getFreshTime")
	}

	token := AccessToken{
		ISS:   issuer,
		TTL:   ttl,
		EXP:   expiry,
		IAT:   issuedAt,
		SUB:   subject,
		Fresh: fresh,
	}

	return token, nil
}

// Parse a refresh token and return a struct with all the claims. Does validation on
// all the claims, including checking if it is expired, has a valid issuer, and
// has the correct scope.
func ParseRefreshToken(config *server.Config, tokenString string) (RefreshToken, error) {
	claims, err := parseToken(config.SecretKey, tokenString)
	if err != nil {
		return RefreshToken{}, errors.Wrap(err, "parseToken")
	}
	expiry, err := checkTokenExpired(claims["exp"])
	if err != nil {
		return RefreshToken{}, errors.Wrap(err, "checkTokenExpired")
	}
	issuer, err := checkTokenIssuer(config.TrustedHost, claims["iss"])
	if err != nil {
		return RefreshToken{}, errors.Wrap(err, "checkTokenIssuer")
	}
	ttl, err := getTokenTTL(claims["ttl"])
	if err != nil {
		return RefreshToken{}, errors.Wrap(err, "getTokenTTL")
	}
	scope, err := getTokenScope(claims["scope"])
	if err != nil {
		return RefreshToken{}, errors.Wrap(err, "getTokenScope")
	}
	if scope != "refresh" {
		return RefreshToken{}, errors.New("Token is not an Refresh token")
	}
	issuedAt, err := getIssuedTime(claims["iat"])
	if err != nil {
		return RefreshToken{}, errors.Wrap(err, "getIssuedTime")
	}
	subject, err := getTokenSubject(claims["sub"])
	if err != nil {
		return RefreshToken{}, errors.Wrap(err, "getTokenSubject")
	}
	jti, err := getTokenJTI(claims["jti"])
	if err != nil {
		return RefreshToken{}, errors.Wrap(err, "getTokenJTI")
	}

	token := RefreshToken{
		ISS: issuer,
		TTL: ttl,
		EXP: expiry,
		IAT: issuedAt,
		SUB: subject,
		JTI: jti,
	}

	return token, nil
}

// Parse a token, validating its signing sigature and returning the claims
func parseToken(secretKey string, tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "jwt.Parse")
	}
	// Token decoded, parse the claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("Failed to parse claims")
	}
	return claims, nil
}

// Check if a token is expired. Returns the expiry if not expired
func checkTokenExpired(expiry interface{}) (int64, error) {
	// Coerce the expiry to a float64 to avoid scientific notation
	expFloat, ok := expiry.(float64)
	if !ok {
		return 0, errors.New("Missing or invalid 'exp' claim")
	}
	// Convert to the int64 time we expect :)
	expiryTime := int64(expFloat)

	// Check if its expired
	isExpired := time.Now().After(time.Unix(expiryTime, 0))
	if isExpired {
		return 0, errors.New("Token has expired")
	}
	return expiryTime, nil
}

// Check if a token has a valid issuer. Returns the issuer if valid
func checkTokenIssuer(trustedHost string, issuer interface{}) (string, error) {
	issuerVal, ok := issuer.(string)
	if !ok {
		return "", errors.New("Missing or invalid 'iss' claim")
	}
	if issuer != trustedHost {
		return "", errors.New("Issuer does not matched trusted host")
	}
	return issuerVal, nil
}

// Check the scope matches the expected scope. Returns scope if true
func getTokenScope(scope interface{}) (string, error) {
	scopeStr, ok := scope.(string)
	if !ok {
		return "", errors.New("Missing or invalid 'scope' claim")
	}
	return scopeStr, nil
}

// Get the TTL of the token, either "session" or "exp"
func getTokenTTL(ttl interface{}) (string, error) {
	ttlStr, ok := ttl.(string)
	if !ok {
		return "", errors.New("Missing or invalid 'ttl' claim")
	}
	if ttlStr != "exp" && ttlStr != "session" {
		return "", errors.New("TTL value is not recognised")
	}
	return ttlStr, nil
}

// Get the time the token was issued at
func getIssuedTime(issued interface{}) (int64, error) {
	// Same float64 -> int64 trick as expiry
	issuedFloat, ok := issued.(float64)
	if !ok {
		return 0, errors.New("Missing or invalid 'iat' claim")
	}
	issuedAt := int64(issuedFloat)
	return issuedAt, nil
}

// Get the freshness expiry timestamp
func getFreshTime(fresh interface{}) (int64, error) {
	freshUntil, ok := fresh.(float64)
	if !ok {
		return 0, errors.New("Missing or invalid 'fresh' claim")
	}
	return int64(freshUntil), nil
}

// Get the subject of the token
func getTokenSubject(sub interface{}) (int, error) {
	subject, ok := sub.(float64)
	if !ok {
		return 0, errors.New("Missing or invalid 'sub' claim")
	}
	return int(subject), nil
}

// Get the JTI of the token
func getTokenJTI(jti interface{}) (uuid.UUID, error) {
	jtiStr, ok := jti.(string)
	if !ok {
		return uuid.UUID{}, errors.New("Missing or invalid 'jti' claim")
	}
	jtiUUID, err := uuid.Parse(jtiStr)
	if err != nil {
		return uuid.UUID{}, errors.New("JTI is not a valid UUID")
	}
	return jtiUUID, nil
}
