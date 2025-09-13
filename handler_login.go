package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jj-attaq/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds *int   `json:"expires_in_seconds,omitempty"`
	}

	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	if err := auth.CheckPasswordHash(params.Password, user.HashedPassword); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	var expirationTime time.Duration
	defaultExpirationTime := time.Duration(time.Hour)

	if params.ExpiresInSeconds == nil || int64(*params.ExpiresInSeconds) > 3600 {
		expirationTime = defaultExpirationTime
	} else {
		expirationTime = time.Duration(int64(*params.ExpiresInSeconds)) * time.Second
	}

	jwtToken, err := auth.MakeJWT(user.ID, cfg.jwtsecret, expirationTime)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error creating JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Token:     jwtToken,
		},
	})
}
