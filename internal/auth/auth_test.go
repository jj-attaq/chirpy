package auth

import (
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

func TestMakeJWT(t *testing.T) {
	secret1 := "testing123"

	tests := []struct {
		name        string
		userID      uuid.UUID
		tokenSecret string
		expiresIn   time.Duration
		wantErr     bool
	}{
		{
			name:        "Successful creation of JWT",
			userID:      uuid.New(),
			tokenSecret: secret1,
			expiresIn:   time.Duration(time.Hour),
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := MakeJWT(tt.userID, tt.tokenSecret, tt.expiresIn)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	secret1 := "testing123"
	// secret2 := "testing456"

	tests := []struct {
		name        string
		userID      uuid.UUID
		tokenString string
		tokenSecret string
		expiresIn   time.Duration
		wantErr     bool
	}{
		{
			name:        "Successful validation of JWT",
			userID:      uuid.New(),
			tokenSecret: secret1,
			expiresIn:   time.Duration(time.Hour),
			wantErr:     false,
		},
	}

	// tokenString, tokenSecret string

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, _ := MakeJWT(tt.userID, tt.tokenSecret, tt.expiresIn)
			tt.tokenString = token
			id, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT error = %v, wantErr %v", err, tt.wantErr)
			}
			if id != tt.userID {
				t.Errorf("tt.userID does not match ValidateJWT output uuid")
			}
		})
	}
}
