package internal

import (
	"context"
	"log"

	pb "racing/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	pb.UnimplementedRacingServer
	pb.UnimplementedStillAliveServer
	db RacingDB
}

func NewServer(conn RacingDB) *Server {
	return &Server{db: conn}
}

func (s *Server) CheckIsRacing(ctx context.Context, in *pb.PlayerMotorcycle) (*pb.RacingStatus, error) {
	return nil, nil
}

func (s *Server) GetHistory(*pb.PlayerUsername, pb.Racing_GetHistoryServer) error {
	return nil
}

func (s *Server) StartMatchmaking(ctx context.Context, in *pb.RaceMotorcycle) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *Server) StillAlive(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	log.Printf("Still Alive")
	return nil, nil
}
