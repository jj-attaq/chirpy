package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jj-attaq/chirpy/internal/database"
)

type Chirp struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    string    `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	// take json body and user id as input
	type parameters struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}

	type response struct {
		Chirp
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// validate said json input
	validChirp, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Chirp", err)
		return
	}
	params.Body = validChirp

	// save the chirp into the database if valid
	id, err := uuid.Parse(params.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Cannot identify user uuid", err)
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   params.Body,
		UserID: id,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't enter Chirp into database", err)
		return
	}
	resp := response{
		Chirp: Chirp{
			ID:        chirp.ID.String(),
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID.String(),
		},
	}

	// respond with 201 and chirp json response
	respondWithJSON(w, http.StatusCreated, resp)
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp too long, must be under 140 characters")
	}

	// Map allows O(1) checking, struct{} as value allows checking with no memory allocation of key value pairs. we only check if the key exists.
	profanity := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	// respondWithJSON(w, http.StatusOK, removeProfanity(params.Body, profanity))
	validChirp := removeProfanity(body, profanity)
	return validChirp, nil
}

func removeProfanity(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		lowerWord := strings.ToLower(word)
		if _, ok := badWords[lowerWord]; ok {
			words[i] = "****"
		}
	}
	result := strings.Join(words, " ")

	return result
}
