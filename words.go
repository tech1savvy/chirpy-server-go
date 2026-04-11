package main

import "strings"

func removeBadWords(chirp string) string {
	badWords := map[string]struct{}{
		"kerfuffle": {}, "sharbert": {}, "fornax": {},
	}

	chirpWords := strings.Split(chirp, " ")
	for i := range chirpWords {
		if _, ok := badWords[strings.ToLower(chirpWords[i])]; ok {
			chirpWords[i] = "****"
		}
	}

	return strings.Join(chirpWords, " ")
}
