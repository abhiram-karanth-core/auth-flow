package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"os"
	"sync"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWK represents a single JSON Web Key
type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// JWKSet is the /.well-known/jwks.json response
type JWKSet struct {
	Keys []JWK `json:"keys"`
}

var (
	privateKey     *rsa.PrivateKey
	keyID          string
	privateKeyOnce sync.Once
)

func loadPrivateKey() {
	privateKeyOnce.Do(func() {
		path := os.Getenv("JWT_PRIVATE_KEY_PATH")
		if path == "" {
			panic("JWT_PRIVATE_KEY_PATH not set")
		}
		keyBytes, err := os.ReadFile(path)
		if err != nil {
			panic("failed to read private key: " + err.Error())
		}
		key, err := jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
		if err != nil {
			panic("failed to parse private key: " + err.Error())
		}
		privateKey = key

		// kid can be set explicitly via env or auto-generated.
		// Explicit kid is useful when rotating — downstream can
		// cache both old and new keys during transition.
		keyID = os.Getenv("JWT_KEY_ID")
		if keyID == "" {
			keyID = uuid.NewString()
		}
	})
}


func GetPublicKey() *rsa.PublicKey {
	return getPrivateKey().Public().(*rsa.PublicKey)
}
func getPrivateKey() *rsa.PrivateKey {
	loadPrivateKey()
	return privateKey
}

func getKeyID() string {
	loadPrivateKey()
	return keyID
}

// JWKSHandler will serve GET /.well-known/jwks.json
// downstream services fetch this once (or on cache miss) to get the public key.
func JWKSHandler(w http.ResponseWriter, r *http.Request) {
	pub := getPrivateKey().Public().(*rsa.PublicKey)

	jwk := JWK{
		Kty: "RSA",
		Use: "sig",
		Kid: getKeyID(),
		Alg: "RS256",
		N:   base64.RawURLEncoding.EncodeToString(pub.N.Bytes()),
		E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(pub.E)).Bytes()),
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600") // downstream can cache for max 1hr
	json.NewEncoder(w).Encode(JWKSet{Keys: []JWK{jwk}})
}
