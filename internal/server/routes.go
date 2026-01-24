package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"authflow/internal/auth"

	authmw "authflow/internal/middleware"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/markbates/goth/gothic"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Get("/auth/{provider}", s.beginAuth)
	r.Get("/auth/{provider}/callback", s.authCallback)
	r.Post("/logout/{provider}", s.logout)
	// protected routes
	r.Group(func(pr chi.Router) {
		pr.Use(authmw.JWTAuth(s.rdb)) 

		pr.Post("/logout", s.logout)

	})
	r.Post("/mint",s.Mint)
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
	tokenStr := extractBearerToken(r)
	if tokenStr == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	token, err := jwt.ParseWithClaims(
		tokenStr,
		&auth.Claims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		},
	)
	if err != nil || !token.Valid {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	claims := token.Claims.(*auth.Claims)

	//revoke JWT
	if err := auth.RevokeToken(s.rdb, claims); err != nil {
		http.Error(w, "logout failed", http.StatusInternalServerError)
		return
	}

	//  clear OAuth session if browser-based
	_ = gothic.Logout(w, r)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("logged out"))
}

func (s *Server) authCallback(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateJWT(user.Email, user.UserID)
	if err != nil {
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	frontendURL := os.Getenv("FRONTEND_URL")

	redirectURL := fmt.Sprintf(
		"%s/oauth/callback?token=%s",
		frontendURL,
		token,
	)

	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

type MintRequest struct {
	Sub      string                 `json:"sub"`
	Provider string                 `json:"provider"`
	Email    string                 `json:"email,omitempty"`
	Claims   map[string]interface{} `json:"claims,omitempty"`
}

func (s *Server) Mint(w http.ResponseWriter, r *http.Request) {
	var req MintRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Sub == "" {
		http.Error(w, "missing sub", http.StatusBadRequest)
		return
	}

	token, err := auth.GenerateJWT(req.Sub, req.Provider)
	if err != nil {
		http.Error(w, "token generation failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": token,
	})
}


