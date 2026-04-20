package main

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsRetreive(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cloudn't retrieve chirps", err)
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt.Time,
			UpdatedAt: dbChirp.UpdatedAt.Time,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerRetreiveChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid chirp ID", err)
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		} else {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirp", err)
		}
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		dbChirp.ID,
		dbChirp.CreatedAt.Time,
		dbChirp.UpdatedAt.Time,
		dbChirp.Body,
		dbChirp.UserID,
	})
}
