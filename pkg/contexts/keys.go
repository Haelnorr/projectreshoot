package contexts

type contextKey string

func (c contextKey) String() string {
	return "projectreshoot context key " + string(c)
}

var (
	contextKeyAuthorizedUser = contextKey("auth-user")
	contextKeyRequestTime    = contextKey("req-time")
)
