package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "jwt-secret!!"
	anotherTokenSecret := "wrong-jwt-secret!!"

	correctTokenString, err := MakeJWT(userID, tokenSecret, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	incorrectTokenString, err := MakeJWT(userID, anotherTokenSecret, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	expiredTokenString, err := MakeJWT(userID, tokenSecret, 1*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		tokenSecret string
		tokenString string
		want        uuid.UUID
		wantErr     bool
	}{
		"valid-token": {
			tokenSecret: tokenSecret,
			tokenString: correctTokenString,
			want:        userID,
			wantErr:     false,
		},
		"invalid-token": {
			tokenSecret: tokenSecret,
			tokenString: incorrectTokenString,
			wantErr:     true,
		},
		"expired-token": {
			tokenSecret: tokenSecret,
			tokenString: expiredTokenString,
			wantErr:     true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ValidateJWT(tc.tokenString, tokenSecret)
			if !tc.wantErr {
				if err != nil {
					t.Fatal(err)
				}
			} else {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			diff := cmp.Diff(tc.want, got)
			if diff != "" {
				t.Error(diff)
			}
		})
	}
}

