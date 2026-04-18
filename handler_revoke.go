package main

import (
	"net/http"

	"github.com/chirpy-server-go/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	// Refresh Token from Header
	headerRefreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	// Revoke Token
	err = cfg.db.RevokeRefreshToken(r.Context(), headerRefreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
