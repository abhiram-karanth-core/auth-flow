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
		var keyBytes []byte

		if b64 := os.Getenv("JWT_PRIVATE_KEY_B64"); b64 != "" {
			// production — key passed as base64 env var
			var err error
			keyBytes, err = base64.StdEncoding.DecodeString(b64)
			if err != nil {
				panic("failed to decode JWT_PRIVATE_KEY_B64: " + err.Error())
			}
		} else if path := os.Getenv("JWT_PRIVATE_KEY_PATH"); path != "" {
			// local dev — key loaded from file
			var err error
			keyBytes, err = os.ReadFile(path)
			if err != nil {
				panic("failed to read private key: " + err.Error())
			}
		} else {
			panic("no private key configured: set JWT_PRIVATE_KEY_B64 or JWT_PRIVATE_KEY_PATH")
		}

		key, err := jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
		if err != nil {
			panic("failed to parse private key: " + err.Error())
		}
		privateKey = key
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
