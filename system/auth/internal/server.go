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
	db *sql.DB // TODO Create dependency
}

func NewServer(conn *sql.DB) *Server {
	return &Server{db: conn}
}

func (s *Server) Login(ctx context.Context, in *pb.PlayerCredentials) (*pb.AuthResult, error) {
	stmt, err := s.db.Prepare("SELECT Username FROM Users WHERE Username=? AND Password=?")
	if err != nil {
		return &pb.AuthResult{Result: false}, err
	}

	var username string
	err = stmt.QueryRow(in.Username, in.Password).Scan(&username)
	if err != nil {
		log.Println("Login result false")
		return &pb.AuthResult{Result: false}, err
	}

	log.Printf("Login result %t", username == in.Username)
	return &pb.AuthResult{Result: username == in.Username}, nil
}

func (s *Server) Register(ctx context.Context, in *pb.PlayerDetails) (*pb.AuthResult, error) {
	stmt, err := s.db.Prepare("INSERT INTO Users VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Println(err)
		return &pb.AuthResult{Result: false}, err
	}

	res, err := stmt.Exec(in.Username, in.Password, in.Email, in.Phone)

	if err != nil {
		log.Println(err)
		log.Printf("Register for username %s result false", in.Username)
		return &pb.AuthResult{Result: false}, err
	}

	rows_affected, err := res.RowsAffected()

	log.Printf("Register for username %s result %t", in.Username, rows_affected != 0)
	return &pb.AuthResult{Result: rows_affected != 0}, err
}

func (s *Server) StillAlive(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	log.Printf("Still Alive")
	return nil, nil
}
