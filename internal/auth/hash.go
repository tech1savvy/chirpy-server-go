package auth

import (
	"errors"

	"github.com/alexedwards/argon2id"
)

func HashPassword(passowrd string) (string, error) {
	hash, err := argon2id.CreateHash(passowrd, argon2id.DefaultParams)
	if err != nil {
		return "", errors.New("failed to hash password")
	}
	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, errors.New("failed to verify password")
	}

	return match, nil
}
