package internal

import (
	"database/sql"

	pb "racing/proto"
)

type Server struct {
	pb.UnimplementedRacingServer
	pb.UnimplementedStillAliveServer
	db *sql.DB // TODO Create dependency
}

func NewServer(conn *sql.DB) *Server {
	return &Server{db: conn}
}
