package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type jwtAuthenticator struct {
	signKey string
}

func NewAuth() Authenticator {
	return &jwtAuthenticator{}
}

type MyCustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (s *jwtAuthenticator) GenerateToken(username string) (string, error) {
	mySigningKey := []byte(s.signKey)

	claims := MyCustomClaims{
		username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(560 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "droplet-api",
			Subject:   "client",
			ID:        uuid.NewString(),
			Audience:  []string{"droplet"},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s *jwtAuthenticator) ParseToken(accessToken string) error {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.signKey), nil
	})
	if err != nil {
		return fmt.Errorf("failed to parse jwt token: %w", err)
	}

	if !token.Valid {
		return fmt.Errorf("token is not valid")
	}

	return nil
}
