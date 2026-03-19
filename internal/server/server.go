package server

import (
	"authflow/ent"
	"authflow/internal/auth"

	"github.com/gorilla/sessions"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	db    *ent.Client
	rdb   *redis.Client
	store *sessions.CookieStore
}

func NewServer(rdb *redis.Client, db *ent.Client) *Server {
	return &Server{
		db:    db,
		rdb:   rdb,
		store: auth.NewAuth(),
	}
}
