package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestAuth(t *testing.T) {
	jwtSecret := []byte("test-secret")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify that we can get claims from context
		claims, err := GetUserFromContext(r.Context())
		if err != nil {
			t.Errorf("Failed to get user from context: %v", err)
			http.Error(w, "Failed to get user from context", http.StatusInternalServerError)
			return
		}
		if claims.Email == "" {
			t.Error("Expected email in claims, got empty string")
		}
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name           string
		token          string
		expectedStatus int
	}{
		{
			name:           "No token provided",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid token format",
			token:          "invalidtoken",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Valid token",
			token:          createValidToken(t, jwtSecret),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Expired token",
			token:          createExpiredToken(t, jwtSecret),
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			rr := httptest.NewRecorder()
			handler := Auth(jwtSecret)(nextHandler)
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tt.expectedStatus)
			}
		})
	}
}

func createValidToken(t *testing.T, secret []byte) string {
	claims := &Claims{
		UserID: "test-user",
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		t.Fatalf("Error creating test token: %v", err)
	}
	return tokenString
}

func createExpiredToken(t *testing.T, secret []byte) string {
	claims := &Claims{
		UserID: "test-user",
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-time.Hour * 2)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-time.Hour * 2)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		t.Fatalf("Error creating expired token: %v", err)
	}
	return tokenString
}

func TestGetUserFromContext(t *testing.T) {
	testClaims := &Claims{
		UserID: "test-user",
		Email:  "test@example.com",
	}

	ctx := context.WithValue(context.Background(), UserContextKey, testClaims)

	claims, err := GetUserFromContext(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if claims.UserID != testClaims.UserID {
		t.Errorf("Expected user ID %s, got %s", testClaims.UserID, claims.UserID)
	}
	if claims.Email != testClaims.Email {
		t.Errorf("Expected email %s, got %s", testClaims.Email, claims.Email)
	}

	_, err = GetUserFromContext(context.Background())
	if err == nil {
		t.Error("Expected error for context without claims, got nil")
	}
}