package internal

import (
	"context"
	"database/sql"
	"log"

	pb "auth/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	pb.UnimplementedAuthenticationServer
	pb.UnimplementedStillAliveServer
	db *sql.DB
}

func NewServer(conn *sql.DB) *Server {
	return &Server{db: conn}
}

func (s *Server) Login(ctx context.Context, in *pb.PlayerCredentials) (*pb.LoginResult, error) {
	stmt, err := s.db.Prepare("SELECT Username FROM Users WHERE Username=? AND Password=?")
	if err != nil {
		log.Println(err)
		return &pb.LoginResult{Result: false}, err
	}

	var username string
	err = stmt.QueryRow(in.Username, in.Password).Scan(username)
	if err != nil {
		log.Println(err)
		return &pb.LoginResult{Result: false}, err
	}

	log.Printf("Login: %t", username == in.Username)

	return &pb.LoginResult{Result: username == in.Username}, nil
}

func (s *Server) Register(ctx context.Context, in *pb.PlayerDetails) (*pb.RegisterResult, error) {
	return &pb.RegisterResult{Result: true, Id: 1}, nil
}

func (s *Server) StillAlive(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	log.Printf("ALIVE")
	return nil, nil
}
