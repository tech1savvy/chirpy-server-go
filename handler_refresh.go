package main

import (
	"net/http"
	"time"

	"github.com/chirpy-server-go/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type returnValues struct {
		Token string `json:"token"`
	}

	headerRefreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), headerRefreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get user for refresh token", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create access token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnValues{
		Token: accessToken,
	})
}
