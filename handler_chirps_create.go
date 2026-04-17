package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/chirpy-server-go/internal/auth"
	"github.com/chirpy-server-go/internal/database"
)

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnValues struct {
		Chirp
	}

	// Authenticate
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	validChirp, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	dbChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   validChirp,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, returnValues{
		Chirp{
			dbChirp.ID,
			dbChirp.CreatedAt.Time,
			dbChirp.UpdatedAt.Time,
			dbChirp.Body,
			dbChirp.UserID,
		},
	})
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	cleaned := removeBadWords(body, badWords)
	return cleaned, nil
}

func removeBadWords(chirp string, badWords map[string]struct{}) string {
	chirpWords := strings.Split(chirp, " ")
	for i := range chirpWords {
		if _, ok := badWords[strings.ToLower(chirpWords[i])]; ok {
			chirpWords[i] = "****"
		}
	}

	return strings.Join(chirpWords, " ")
}
