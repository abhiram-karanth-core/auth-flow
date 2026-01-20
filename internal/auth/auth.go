package auth

import (
	"os"

	"github.com/joho/godotenv"

	"net/http"

	"github.com/gorilla/sessions"
)

const (
	MaxAge = 86400 * 30
)

func NewAuth() *sessions.CookieStore {
	_ = godotenv.Load()


	secretKey := os.Getenv("SESSION_SECRET")
	isProd := os.Getenv("ENV") == "production"

	store := sessions.NewCookieStore([]byte(secretKey))
	store.MaxAge(MaxAge)

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   MaxAge,
		HttpOnly: true,
		Secure:   isProd,
		SameSite: http.SameSiteLaxMode,
	}
	return store
}
