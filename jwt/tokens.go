package jwt

import "github.com/google/uuid"

type Token interface {
	GetJTI() uuid.UUID
	GetEXP() int64
	GetScope() string
}

// Access token
type AccessToken struct {
	ISS   string    // Issuer, generally TrustedHost
	IAT   int64     // Time issued at
	EXP   int64     // Time expiring at
	TTL   string    // Time-to-live: "session" or "exp". Used with 'remember me'
	SUB   int       // Subject (user) ID
	JTI   uuid.UUID // UUID-4 used for identifying blacklisted tokens
	Fresh int64     // Time freshness expiring at
	Scope string    // Should be "access"
}

// Refresh token
type RefreshToken struct {
	ISS   string    // Issuer, generally TrustedHost
	IAT   int64     // Time issued at
	EXP   int64     // Time expiring at
	TTL   string    // Time-to-live: "session" or "exp". Used with 'remember me'
	SUB   int       // Subject (user) ID
	JTI   uuid.UUID // UUID-4 used for identifying blacklisted tokens
	Scope string    // Should be "refresh"
}

func (a AccessToken) GetJTI() uuid.UUID {
	return a.JTI
}
func (r RefreshToken) GetJTI() uuid.UUID {
	return r.JTI
}
func (a AccessToken) GetEXP() int64 {
	return a.EXP
}
func (r RefreshToken) GetEXP() int64 {
	return r.EXP
}
func (a AccessToken) GetScope() string {
	return a.Scope
}
func (r RefreshToken) GetScope() string {
	return r.Scope
}
