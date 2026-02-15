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
	r.Post("/mint", s.Mint)
	return r
}
func extractOAuthParams(r *http.Request) (string, string, error) {
	clientID := r.URL.Query().Get("client_id")
	redirectURI := r.URL.Query().Get("redirect_uri")

	if clientID == "" || redirectURI == "" {
		return "", "", fmt.Errorf("missing client_id or redirect_uri")
	}
	return clientID, redirectURI, nil
}
func (s *Server) isAllowedRedirect(clientID, redirectURI string) bool {
	allowed := map[string][]string{
		"ragworks": {
			"https://ragworks-wheat.vercel.app/oauth/callback",
			"http://localhost:3000/oauth/callback",
		},
	}

	uris, ok := allowed[clientID]
	if !ok {
		return false
	}

	for _, uri := range uris {
		if uri == redirectURI {
			return true
		}
	}
	return false
}

func (s *Server) beginAuth(w http.ResponseWriter, r *http.Request) {
	clientID, redirectURI, err := extractOAuthParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !s.isAllowedRedirect(clientID, redirectURI) {
		http.Error(w, "invalid redirect_uri", http.StatusBadRequest)
		return
	}
	session, _ := s.store.Get(r, "authflow")
	session.Values["client_id"] = clientID
	session.Values["redirect_uri"] = redirectURI
	session.Save(r, w)

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

	session, _ := s.store.Get(r, "authflow")
	redirectURI, ok := session.Values["redirect_uri"].(string)
	if !ok {
		http.Error(w, "redirect uri missing", http.StatusBadRequest)
		return
	}
	redirectURL := fmt.Sprintf("%s?token=%s",redirectURI,token)

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

	token, err := auth.GenerateJWT(req.Sub, req.Sub)
	if err != nil {
		http.Error(w, "token generation failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": token,
	})
}
