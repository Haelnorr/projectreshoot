package jwt

import "github.com/google/uuid"

// Access token
type AccessToken struct {
	ISS   string
	Scope string
	IAT   int64
	EXP   int64
	SUB   int
	Fresh int64
	Roles []string
}

// Refresh token
type RefreshToken struct {
	ISS   string
	Scope string
	IAT   int64
	EXP   int64
	SUB   int
	JTI   uuid.UUID
}
