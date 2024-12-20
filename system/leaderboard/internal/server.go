package internal

import (
	"context"
	"errors"
	"log"

	pb "leaderboard/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

type LeaderboardInfo struct {
	username string
	points   int32
	position int32
}

type Server struct {
	pb.UnimplementedLeaderboardServer
	pb.UnimplementedStillAliveServer
	db LeaderboardDB
}

func NewServer(conn LeaderboardDB) *Server {
	return &Server{db: conn}
}

func (s *Server) GetLeaderboard(_ *emptypb.Empty, stream pb.Leaderboard_GetLeaderboardServer) error {
	leaderboard := s.db.GetLeaderboard()

	if len(leaderboard) == 0 || leaderboard == nil {
		return nil
	}

	for i := 0; i < len(leaderboard); i++ {
		stream.Send(&pb.LeaderboardPosition{
			Username: leaderboard[i].username,
			Position: leaderboard[i].position,
			Points:   leaderboard[i].points,
		})
	}

	return errors.New("no user in leaderboard")
}

func (s *Server) GetPlayer(ctx context.Context, in *pb.PlayerUsername) (*pb.LeaderboardPosition, error) {
	user, err := s.db.GetUserInfo(in.Username)

	if err != nil || user == nil {
		return nil, err
	}

	return &pb.LeaderboardPosition{Username: user.username, Position: user.position, Points: user.points}, nil
}

func (s *Server) AddPoints(ctx context.Context, in *pb.PointIncrement) (*emptypb.Empty, error) {
	return nil, s.db.IncrementPoints(in.Username, int(in.Points))
}

func (s *Server) StillAlive(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	log.Printf("Still Alive")
	return nil, nil
}
