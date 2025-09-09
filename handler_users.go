package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type emailReq struct {
	Email string `json:"email"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	email := emailReq{}

	err := decoder.Decode(&email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	// user requests context for timeout
	user, err := cfg.db.CreateUser(r.Context(), email.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
