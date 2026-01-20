package main

import (
	"authflow/internal/auth"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"authflow/internal/server"
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
			"http://localhost:8080/auth/google/callback",
		),
	)
	srv := server.NewServer() 
	handler := srv.RegisterRoutes()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
