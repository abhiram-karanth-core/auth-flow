package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	Email    string `json:"email"`
	Provider string `json:"provider"`
	jwt.RegisteredClaims
}

func GenerateJWT(email, providerUserID string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET not set")
	}
	// now := time.Now()
	jti := uuid.NewString()
	claims := Claims{
		Email:    email,
		Provider: "google",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   providerUserID,
			ID:        jti,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "authflow-go",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
