package main

import (
	"encoding/json"
	"net/http"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		Valid bool `json:"valid"`
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

	if validitiy {
		respondWithJSON(w, http.StatusOK, returnVals{
			Valid: true,
		})
	} else {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
	}
}
