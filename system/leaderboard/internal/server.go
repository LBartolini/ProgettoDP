package internal

import (
	"database/sql"

	pb "leaderboard/proto"
)

type Server struct {
	pb.UnimplementedLeaderboardServer
	pb.UnimplementedStillAliveServer
	db *sql.DB // TODO Create dependency
}

func NewServer(conn *sql.DB) *Server {
	return &Server{db: conn}
}
