package middleware

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync/atomic"
	"testing"

	"projectreshoot/db"
	"projectreshoot/tests"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPageLoginRequired(t *testing.T) {
	cfg, err := tests.TestConfig()
	require.NoError(t, err)
	logger := tests.NilLogger()
	ver, err := strconv.ParseInt(cfg.DBName, 10, 0)
	require.NoError(t, err)
	conn, err := tests.SetupTestDB(ver)
	require.NoError(t, err)
	sconn := db.MakeSafe(conn, logger)
	defer sconn.Close()

	// Handler to check outcome of Authentication middleware
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	var maint uint32
	atomic.StoreUint32(&maint, 0)
	// Add the middleware and create the server
	loginRequiredHandler := LoginReq(testHandler)
	authHandler := Authentication(logger, cfg, sconn, loginRequiredHandler, &maint)
	server := httptest.NewServer(authHandler)
	defer server.Close()

	tokens := getTokens()

	tests := []struct {
		name         string
		accessToken  string
		refreshToken string
		expectedCode int
	}{
		{
			name:         "Valid Login",
			accessToken:  tokens["accessFresh"],
			refreshToken: "",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Expired login",
			accessToken:  tokens["accessExpired"],
			refreshToken: tokens["refreshExpired"],
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "No login",
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
		})
	}
}
