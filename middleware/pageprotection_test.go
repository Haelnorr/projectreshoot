package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"projectreshoot/tests"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPageLoginRequired(t *testing.T) {
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
		w.WriteHeader(http.StatusOK)
	})

	// Add the middleware and create the server
	loginRequiredHandler := RequiresLogin(testHandler)
	authHandler := Authentication(logger, cfg, conn, loginRequiredHandler)
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
