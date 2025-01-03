package internal

import (
	"context"

	pb "auth/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	pb.UnimplementedAuthenticationServer
	pb.UnimplementedStillAliveServer
	db AuthDB
}

func NewServer(conn AuthDB) *Server {
	return &Server{db: conn}
}

func (s *Server) Login(ctx context.Context, in *pb.PlayerCredentials) (*pb.AuthResult, error) {
	res, err := s.db.Login(in.Username, in.Password)

	return &pb.AuthResult{Result: res}, err
}

func (s *Server) Register(ctx context.Context, in *pb.PlayerDetails) (*pb.AuthResult, error) {
	res, err := s.db.Register(in.Username, in.Password, in.Email, in.Phone)

	return &pb.AuthResult{Result: res}, err
}

func (s *Server) StillAlive(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, nil
}
