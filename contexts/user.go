package contexts

import (
	"context"
	"projectreshoot/db"
)

// Return a new context with the user added in
func SetUser(ctx context.Context, u *db.User) context.Context {
	return context.WithValue(ctx, contextKeyAuthorizedUser, u)
}

// Retrieve a user from the given context. Returns nil if not set
func GetUser(ctx context.Context) *db.User {
	user, ok := ctx.Value(contextKeyAuthorizedUser).(*db.User)
	if !ok {
		return nil
	}
	return user
}
