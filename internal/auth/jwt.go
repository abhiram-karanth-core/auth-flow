package auth

import (

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

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = getKeyID() // downstream uses kid to pick the right key from JWKS
	return token.SignedString(getPrivateKey())
}
