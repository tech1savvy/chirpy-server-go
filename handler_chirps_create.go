package main

import (
	"encoding/json"
	"net/http"

	"github.com/chirpy-server-go/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	validitiy := len(params.Body) <= maxChirpLength

	if !validitiy {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
	}

	dbChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   params.Body,
		UserID: params.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}
	chirp := Chirp{
		dbChirp.ID,
		dbChirp.CreatedAt.Time,
		dbChirp.UpdatedAt.Time,
		dbChirp.Body,
		dbChirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}
