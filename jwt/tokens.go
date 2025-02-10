package jwt

import "github.com/google/uuid"

// Access token
type AccessToken struct {
	ISS   string // Issuer, generally TrustedHost
	IAT   int64  // Time issued at
	EXP   int64  // Time expiring at
	TTL   string // Time-to-live: "session" or "exp". Used with 'remember me'
	SUB   int    // Subject (user) ID
	Fresh int64  // Time freshness expiring at
}

// Refresh token
type RefreshToken struct {
	ISS string    // Issuer, generally TrustedHost
	IAT int64     // Time issued at
	EXP int64     // Time expiring at
	TTL string    // Time-to-live: "session" or "exp". Used with 'remember me'
	SUB int       // Subject (user) ID
	JTI uuid.UUID // UUID-4 used for identifying blacklisted refresh tokens
}
