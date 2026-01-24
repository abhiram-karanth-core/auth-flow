package server

import (
	"authflow/internal/database"

	"github.com/redis/go-redis/v9"
)

type Server struct {
	db  database.Service
	rdb *redis.Client
}

func NewServer(rdb *redis.Client) *Server {
	return &Server{
		db:  database.New(),
		rdb: rdb,
	}
}
