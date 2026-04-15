package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRemoveBadWords(t *testing.T) {
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	tests := map[string]struct {
		chirp string
		want  string
	}{
		"without-puntuations": {
			"kerfuffle yo sharbert yo",
			"**** yo **** yo",
		},
		"with-puntuations": {
			"kerfuffle! yo sharbert yo",
			"kerfuffle! yo **** yo",
		},
		"with-capitals": {
			"kerfuffle Yo sharbert yo",
			"**** Yo **** yo",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := removeBadWords(tc.chirp, badWords)
			diff := cmp.Diff(tc.want, got)
			if diff != "" {
				t.Error(diff)
			}
		})
	}
}
