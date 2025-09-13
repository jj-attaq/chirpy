package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret", time.Hour)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	userID := uuid.New()
	token, _ := MakeJWT(userID, "secret", time.Hour)

	validHeader := http.Header{"Authorization": []string{"Bearer " + token}}

	tests := []struct {
		name    string
		header  http.Header
		wantErr bool
	}{
		{
			name:    "Success",
			header:  validHeader,
			wantErr: false,
		},
		{
			name:    "No authorization header",
			header:  http.Header{"NoAuthorizationHeader": []string{"Bearer " + token}},
			wantErr: true,
		},
		{
			name:    "Empty authorization header",
			header:  http.Header{"Authorization": []string{}},
			wantErr: true,
		},
		{
			name: "Multiple authorization headers",
			header: http.Header{"Authorization": []string{
				"Bearer " + token,
				"extra element in slice"},
			},
			wantErr: true,
		},
		{
			name:    "Invalid prefix",
			header:  http.Header{"Authorization": []string{"Invalid prefix" + token}},
			wantErr: true,
		},
		{
			name:    "Empty bearer token",
			header:  http.Header{"Authorization": []string{"Bearer " + ""}},
			wantErr: true,
		},
		{
			name:    "Whitespace around value is trimmed",
			header:  http.Header{"Authorization": []string{"   Bearer " + token + "   "}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		_, err := GetBearerToken(tt.header)
		if (err != nil) != tt.wantErr {
			t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
}
