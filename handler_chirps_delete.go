package main

import (
	"net/http"

	"github.com/chirpy-server-go/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleChirpsDelete(w http.ResponseWriter, r *http.Request) {
	// Query Paramters
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cloudn't parse chirpID", err)
		return
	}

	// Validate JWT
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Cloudn't get JWT", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Cloudn't validate JWT", err)
		return
	}

	// Get Chirp
	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Cloudn't find chirp", err)
		return
	}
	// Validate Author
	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You are not the author", err)
		return
	}

	// Delete Chirp from DB
	err = cfg.db.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cloudn't delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
