package services

import (
	"context"
	pb "orchestrator/proto"
	"time"

	"google.golang.org/grpc"
)

type Auth interface {
	StillAlive
	Login(username string, password string) (bool, error)
	Register(username, password, email, phone string) (bool, error)
}

// gRPC implementation of Auth interface
type AuthService struct {
	conn *grpc.ClientConn
}

func NewAuthService(conn *grpc.ClientConn) *AuthService {
	return &AuthService{conn: conn}
}

func (s *AuthService) StillAlive() bool {
	return StillAliveHandle(s.conn)
}

func (s *AuthService) Close() {
	s.conn.Close()
}

func (s *AuthService) Login(username string, password string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := pb.NewAuthenticationClient(s.conn).Login(ctx, &pb.PlayerCredentials{Username: username, Password: password})
	if err != nil || res == nil {
		return false, err
	}

	return res.Result, nil
}

func (s *AuthService) Register(username, password, email, phone string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := pb.NewAuthenticationClient(s.conn).Register(ctx, &pb.PlayerDetails{Username: username, Password: password, Email: email, Phone: phone})
	if err != nil || res == nil {
		return false, err
	}

	return res.Result, nil
}
