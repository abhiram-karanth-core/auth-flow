package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"authflow/internal/auth"

	authmw "authflow/internal/middleware"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/markbates/goth/gothic"
	"golang.org/x/crypto/bcrypt"
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
	r.Get("/.well-known/jwks.json", auth.JWKSHandler)
	r.Get("/auth/{provider}", s.beginAuth)
	r.Get("/auth/{provider}/callback", s.authCallback)
	r.Post("/logout/{provider}", s.logout)
	// protected routes
	r.Group(func(pr chi.Router) {
		pr.Use(authmw.JWTAuth(s.rdb))

		pr.Post("/logout", s.logout)

	})
	r.Post("/mint", s.Mint)
	r.Post("/clients", s.RegisterClients)
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

func (s *Server) beginAuth(w http.ResponseWriter, r *http.Request) {
	clientID, redirectURI, err := extractOAuthParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	client, err := s.db.OAuthClient.Get(r.Context(), clientID)
	if err != nil || client.RedirectURI != redirectURI {
		http.Error(w, "invalid client_id or redirect_uri", http.StatusBadRequest)
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
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected alg: %v", t.Header["alg"])
			}
			return auth.GetPublicKey(), nil // derive from private key
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
	redirectURL := fmt.Sprintf("%s?token=%s", redirectURI, token)

	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

type MintRequest struct {
    ClientID     string                 `json:"client_id"`
    ClientSecret string                 `json:"client_secret"`
    Sub          string                 `json:"sub"`
    Provider     string                 `json:"provider"`
    Email        string                 `json:"email,omitempty"`
    Claims       map[string]interface{} `json:"claims,omitempty"`
}

func (s *Server) Mint(w http.ResponseWriter, r *http.Request) {
    var req MintRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }

    if req.Sub == "" || req.ClientID == "" || req.ClientSecret == "" {
        http.Error(w, "missing required fields", http.StatusBadRequest)
        return
    }

    // Look up the registered client
    client, err := s.db.OAuthClient.Get(r.Context(), req.ClientID)
    if err != nil {
        http.Error(w, "invalid client", http.StatusUnauthorized)
        return
    }

    // Verify the secret against the bcrypt hash you already store
    if err := bcrypt.CompareHashAndPassword([]byte(client.Secret), []byte(req.ClientSecret)); err != nil {
        http.Error(w, "invalid client credentials", http.StatusUnauthorized)
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

func (s *Server) RegisterClients(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name        string `json:"name"`
		RedirectUri string `json:"redirect_uri"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	clientID := uuid.NewString()
	secret := uuid.NewString()
	hashed, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	_, err = s.db.OAuthClient.Create().
		SetID(clientID).
		SetName(body.Name).
		SetSecret(string(hashed)).
		SetRedirectURI(body.RedirectUri).
		Save(r.Context())
	if err != nil {
		http.Error(w, "failed to register client", http.StatusInternalServerError)
		return
	}

	// return raw secret only once — not stored in plaintext
	json.NewEncoder(w).Encode(map[string]string{
		"client_id":     clientID,
		"client_secret": secret,
		"redirect_uri":  body.RedirectUri,
	})
}
