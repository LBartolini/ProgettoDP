package services

import (
	"context"
	"io"
	"log"
	"time"

	pb "orchestrator/proto"

	"google.golang.org/grpc"
)

type LeaderboardPosition struct {
	Username string
	Points   int
	Position int
}

type Leaderboard interface {
	StillAlive
	AddPoints(username string, points int) error
	GetPlayer(username string) (*LeaderboardPosition, error)
	GetFullLeaderboard() ([]*LeaderboardPosition, error)
}

// gRPC implementation of Leaderboard interface
type LeaderboardService struct {
	conn *grpc.ClientConn
}

func NewLeaderboardService(conn *grpc.ClientConn) *LeaderboardService {
	return &LeaderboardService{conn: conn}
}

func (s *LeaderboardService) StillAlive() bool {
	return StillAliveHandle(s.conn)
}

func (s *LeaderboardService) Close() {
	s.conn.Close()
}

func (s *LeaderboardService) GetFullLeaderboard() ([]*LeaderboardPosition, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := pb.NewLeaderboardClient(s.conn).GetFullLeaderboard(ctx, nil)
	if err != nil {
		return nil, err
	}

	var leaderboard []*LeaderboardPosition
	for {
		p, err := r.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println(err)
			return nil, err
		} else {
			pos := &LeaderboardPosition{Username: p.Username, Points: int(p.Points), Position: int(p.Position)}
			leaderboard = append(leaderboard, pos)
		}
	}

	return leaderboard, nil
}

func (s *LeaderboardService) GetPlayer(username string) (*LeaderboardPosition, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	pos, err := pb.NewLeaderboardClient(s.conn).GetPlayer(ctx, &pb.PlayerUsername{Username: username})
	if err != nil {
		return nil, err
	}
	return &LeaderboardPosition{Username: pos.Username, Points: int(pos.Points), Position: int(pos.Position)}, nil
}

func (s *LeaderboardService) AddPoints(username string, points int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := pb.NewLeaderboardClient(s.conn).AddPoints(ctx, &pb.PointIncrement{Username: username, Points: int32(points)})
	return err
}
