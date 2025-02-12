package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"projectreshoot/contexts"
	"projectreshoot/db"
	"projectreshoot/jwt"
	"projectreshoot/tests"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthenticationMiddleware(t *testing.T) {
	// Basic setup
	cfg, err := tests.TestConfig()
	require.NoError(t, err)
	logger := tests.NilLogger()
	conn, err := tests.SetupTestDB()
	require.NoError(t, err)
	require.NotNil(t, conn)
	defer tests.DeleteTestDB()

	// Handler to check outcome of Authentication middleware
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := contexts.GetUser(r.Context())
		if user == nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(strconv.Itoa(0)))
			return
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(strconv.Itoa(user.ID)))
		}
	})

	// Add the middleware and create the server
	authHandler := Authentication(logger, cfg, conn, testHandler)
	require.NoError(t, err)
	server := httptest.NewServer(authHandler)
	defer server.Close()

	// Setup the user and tokens to test with
	user, err := db.GetUserFromID(conn, 1)
	require.NoError(t, err)

	// Good tokens
	atStr, _, err := jwt.GenerateAccessToken(cfg, &user, false, false)
	require.NoError(t, err)
	rtStr, _, err := jwt.GenerateRefreshToken(cfg, &user, false)
	require.NoError(t, err)

	// Create a token and revoke it for testing
	expStr, _, err := jwt.GenerateAccessToken(cfg, &user, false, false)
	require.NoError(t, err)
	expT, err := jwt.ParseAccessToken(cfg, conn, expStr)
	require.NoError(t, err)
	err = jwt.RevokeToken(conn, expT)
	require.NoError(t, err)

	// Make sure it actually got revoked
	expT, err = jwt.ParseAccessToken(cfg, conn, expStr)
	require.Error(t, err)

	tests := []struct {
		name         string
		id           int
		accessToken  string
		refreshToken string
		expectedCode int
	}{
		{
			name:         "Valid Access Token",
			id:           1,
			accessToken:  atStr,
			refreshToken: "",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Valid Refresh Token (Triggers Refresh)",
			id:           1,
			accessToken:  expStr,
			refreshToken: rtStr,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Refresh token revoked (after refresh)",
			accessToken:  expStr,
			refreshToken: rtStr,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "Invalid Tokens",
			accessToken:  expStr,
			refreshToken: expStr,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "No Tokens",
			accessToken:  "",
			refreshToken: "",
			expectedCode: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &http.Client{}

			req, _ := http.NewRequest(http.MethodGet, server.URL, nil)

			// Add cookies if provided
			if tt.accessToken != "" {
				req.AddCookie(&http.Cookie{Name: "access", Value: tt.accessToken})
			}
			if tt.refreshToken != "" {
				req.AddCookie(&http.Cookie{Name: "refresh", Value: tt.refreshToken})
			}

			resp, err := client.Do(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.StatusCode)
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Equal(t, strconv.Itoa(tt.id), string(body))
		})
	}
}
