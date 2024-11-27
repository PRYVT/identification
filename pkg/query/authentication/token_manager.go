package authentication

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type TokenManager struct {
	signingSecret string
}

func NewTokenManager() (*TokenManager, error) {
	secret := os.Getenv("SIGNING_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("SIGNING_SECRET environment variable not set")
	}
	return &TokenManager{signingSecret: secret}, nil
}

func (tm *TokenManager) CreateToken(userUuid uuid.UUID) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userUuid,
		"iss": "pryvt",
		"aud": "local-audience",
		"exp": time.Now().Add(time.Second * 30).Unix(),
		"iat": time.Now().Unix(),
	})

	signedKey, err := t.SignedString([]byte(tm.signingSecret))
	if err != nil {
		log.Err(err).Msg("Error while signing token")
		return "", fmt.Errorf("error while signing token")
	}
	return signedKey, nil
}
func (tm *TokenManager) VerifyToken(tokenString string) (*jwt.Token, error) {
	// Parse the token with the secret key

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(tm.signingSecret), nil
	})

	// Check for verification errors
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Return the verified token
	return token, nil
}
