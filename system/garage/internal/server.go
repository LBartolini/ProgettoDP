package internal

import (
	"database/sql"

	pb "garage/proto"
)

type Server struct {
	pb.UnimplementedGarageServer
	pb.UnimplementedStillAliveServer
	db *sql.DB // TODO Create dependency
}

func NewServer(conn *sql.DB) *Server {
	return &Server{db: conn}
}
