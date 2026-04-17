package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/chirpy-server-go/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"Password"`
		ExpiresInSeconds int    `json:"expires_in_seconds,omitempty"`
	}

	type returnValues struct {
		User
		Token         string `json:"token"`
		RefereshToken string `json:"refresh_token"`
	}

	// Decode Params
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cloudn't decode parameters", err)
		return
	}

	expirationTime := time.Hour
	if params.ExpiresInSeconds > 0 && params.ExpiresInSeconds < 60*60 {
		expirationTime = time.Duration(params.ExpiresInSeconds) * time.Second
	}

	// Authenticate
	dbUser, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	isCorrect, err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	if !isCorrect {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	// JWT
	accessToken, err := auth.MakeJWT(dbUser.ID, cfg.jwtSecret, expirationTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create access token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnValues{
		User: User{
			ID:        dbUser.ID,
			CreatedAt: dbUser.CreatedAt.Time,
			UpdateAt:  dbUser.UpdatedAt.Time,
			Email:     dbUser.Email,
		},
		Token: accessToken,
	})
}
