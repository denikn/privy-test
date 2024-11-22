package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v4"
)

func TestJWTMiddleware(t *testing.T) {
	// MockClaims is a mock implementation of jwt.Claims for testing purposes
	type MockClaims struct {
		Username string `json:"username"`
		jwt.StandardClaims
	}

	// MockHandlerFunc is a mock implementation of http.HandlerFunc for testing purposes
	mockHandlerFunc := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	// Set JWT_SECRET_KEY for testing environment
	os.Setenv("JWT_SECRET_KEY", "test_secret_key")
	defer os.Unsetenv("JWT_SECRET_KEY")

	tests := []struct {
		name           string
		authorization  string
		expectedStatus int
	}{
		{
			name:           "Missing authorization header",
			authorization:  "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Malformed token",
			authorization:  "invalid_token",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid token",
			authorization:  "Bearer invalid_token",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", tt.authorization)
			rec := httptest.NewRecorder()

			// Call JWTMiddleware with mockHandlerFunc as next handler
			JWTMiddleware(http.HandlerFunc(mockHandlerFunc)).ServeHTTP(rec, req)

			// Check response status code
			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d; got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}
