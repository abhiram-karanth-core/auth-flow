package server

import (
	"fmt"
	"net/http"
	"os"

	"authflow/internal/auth"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/markbates/goth/gothic"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Get("/auth/{provider}", s.beginAuth)
	r.Get("/auth/{provider}/callback", s.authCallback)
	r.Get("/logout/{provider}", s.logout)
	return r
}
func (s *Server) beginAuth(w http.ResponseWriter, r *http.Request) {
	// Try to complete auth (user already logged in)
	if _, err := gothic.CompleteUserAuth(w, r); err == nil {
		w.Write([]byte("Already authenticated"))
		return
	}

	// Start OAuth flow
	gothic.BeginAuthHandler(w, r)
}
func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	_ = gothic.Logout(w, r)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (s *Server) authCallback(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateJWT(user.Email)
	if err != nil {
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	// hardcoded frontend url as of now
	frontendURL := os.Getenv("FRONTEND_URL")

	redirectURL := fmt.Sprintf(
		"%s/oauth/callback?token=%s&username=%s",
		frontendURL,
		token,
		user.Email,
	)

	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}
