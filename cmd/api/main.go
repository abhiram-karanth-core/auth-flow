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

	err := rdb.Set(redisclient.Ctx, "session:user:123", "logged_in", time.Minute*10).Err()
	if err != nil {
		log.Fatal(err)
	}

	srv := server.NewServer(rdb)
	handler := srv.RegisterRoutes()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
