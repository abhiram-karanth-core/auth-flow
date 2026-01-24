package main

import (
	"authflow/internal/auth"
	redisclient "authflow/internal/redis"
	"log"
	"net/http"
	"os"
	"time"

	"authflow/internal/server"

	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func main() {
	store := auth.NewAuth()
	gothic.Store = store
	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	gothic.GetProviderName = func(r *http.Request) (string, error) {
		return chi.URLParam(r, "provider"), nil
	}
	goth.UseProviders(
		google.New(
			googleClientId,
			googleClientSecret,
			"https://authflow-go.onrender.com/auth/google/callback",
		),
	)
	rdb := redisclient.New()

	err := rdb.Set(redisclient.Ctx, "redis:test", "it_works", 5*time.Minute).Err()
	if err != nil {
		log.Fatal("Redis SET failed:", err)
	}

	val, err := rdb.Get(redisclient.Ctx, "redis:test").Result()
	if err != nil {
		log.Fatal("Redis GET failed:", err)
	}

	log.Println("Redis test value:", val)

	srv := server.NewServer(rdb)
	handler := srv.RegisterRoutes()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
