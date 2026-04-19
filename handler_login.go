package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/chirpy-server-go/internal/auth"
	"github.com/chirpy-server-go/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"Password"`
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
	accessToken, err := auth.MakeJWT(dbUser.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create access token", err)
		return
	}

	// Referesh Token
	refreshToken := auth.MakeRefreshToken()
	err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    dbUser.ID,
		ExpiresAt: time.Now().UTC().Add(60 * 24 * time.Hour),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to save refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnValues{
		User: User{
			ID:          dbUser.ID,
			CreatedAt:   dbUser.CreatedAt.Time,
			UpdateAt:    dbUser.UpdatedAt.Time,
			Email:       dbUser.Email,
			IsChirpyRed: dbUser.IsChirpyRed,
		},
		Token:         accessToken,
		RefereshToken: refreshToken,
	})
}
