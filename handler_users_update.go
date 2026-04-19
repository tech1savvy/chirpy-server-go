package main

import (
	"encoding/json"
	"net/http"

	"github.com/chirpy-server-go/internal/auth"
	"github.com/chirpy-server-go/internal/database"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type returnVaules struct {
		User
	}

	// Auth
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Cloudn't find JWT", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Cloudn't validate JWT", err)
		return
	}

	// Decode Parameters
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cloudn't decode parameters", err)
		return
	}

	// Hash New Password
	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cloudn't hash password", err)
		return
	}

	// Update User in DB
	user, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hash,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to update user resource", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVaules{
		User{
			ID:          user.ID,
			Email:       user.Email,
			CreatedAt:   user.CreatedAt.Time,
			UpdateAt:    user.UpdatedAt.Time,
			IsChirpyRed: user.IsChirpyRed,
		},
	})
}
