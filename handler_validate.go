package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

var profanity = []string{"kerfuffle", "sharbert", "fornax"}

type parameters struct {
	Body string `json:"body"`
}

type validBody struct {
	Valid bool `json:"valid"`
}

type cleanParams struct {
	Body string `json:"cleaned_body"`
}

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, cleanParams{
		Body: params.removeProfanity().Body,
	})

	// respondWithJSON(w, http.StatusOK, validBody{
	// 	Valid: true,
	// })
}

func (p *parameters) removeProfanity() *parameters {
	words := strings.Split(p.Body, " ")
	for i, word := range words {
		lowerWord := strings.ToLower(word)
		for _, profaneWord := range profanity {
			if lowerWord == profaneWord {
				words[i] = "****"
			}
		}
	}
	newParams := strings.Join(words, " ")
	p.Body = newParams

	return p
}

// func (p *parameters) removeProfanity() *parameters {
// 	lowerP := strings.ToLower(p.Body)
// 	for _, word := range profanity {
// 		if strings.Contains(lowerP, word) {
// 			strings.ReplaceAll(lowerP, word, "****")
// 		}
// 	}
// 	p.Body = lowerP
//
// 	return p
// }
