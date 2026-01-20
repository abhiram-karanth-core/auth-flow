package server

import "authflow/internal/database"

type Server struct {
	db database.Service
}

func NewServer() *Server {
	return &Server{
		db: database.New(),
	}
}
