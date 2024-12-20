package internal

import (
	"context"
	"database/sql"
	"log"

	pb "racing/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	pb.UnimplementedRacingServer
	pb.UnimplementedStillAliveServer
	db *sql.DB // TODO Create dependency
}

func NewServer(conn *sql.DB) *Server {
	return &Server{db: conn}
}

func (s *Server) StillAlive(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	log.Printf("Still Alive")
	return nil, nil
}
