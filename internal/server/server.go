package server

import (
	"authflow/internal/database"

	"github.com/gorilla/sessions"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	db    database.Service
	rdb   *redis.Client
	store *sessions.CookieStore
}

func NewServer(rdb *redis.Client) *Server {
	return &Server{
		db:  database.New(),
		rdb: rdb,
	}
}
