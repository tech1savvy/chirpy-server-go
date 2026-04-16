package auth

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestHashPassword(t *testing.T) {
	tests := map[string]struct {
		password string
	}{
		"succesful hash generation": {
			password: "pa$$word",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := HashPassword(tc.password)
			if err != nil {
				t.Fail()
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "pa$$word"
	invalidPassword := "wrong"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		password       string
		hashedPassword string
		want           bool
	}{
		"valid password": {
			password,
			hash,
			true,
		},
		"invalid password": {
			invalidPassword,
			hash,
			false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := CheckPasswordHash(tc.password, tc.hashedPassword)
			if err != nil {
				t.Fail()
			}
			diff := cmp.Diff(tc.want, got)
			if diff != "" {
				t.Error(diff)
			}
		})
	}
}
