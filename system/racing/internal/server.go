package internal

import (
	"context"
	"log"
	"time"

	pb "racing/proto"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedRacingServer
	pb.UnimplementedStillAliveServer
	db           RacingDB
	orchestrator *grpc.ClientConn
}

func NewServer(conn RacingDB, orchestrator *grpc.ClientConn) *Server {
	return &Server{db: conn, orchestrator: orchestrator}
}

func (s *Server) CheckIsRacing(ctx context.Context, in *pb.PlayerMotorcycle) (*pb.RacingStatus, error) {
	track, _ := s.db.CheckIsRacing(in.Username, int(in.MotorcycleId))

	log.Printf("Checking racing (%s:%d), status %t", in.Username, in.MotorcycleId, track != "")

	return &pb.RacingStatus{IsRacing: track != "", TrackName: track}, nil
}

func (s *Server) GetHistory(in *pb.PlayerUsername, stream pb.Racing_GetHistoryServer) error {
	results, err := s.db.GetHistory(in.Username)

	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Retrieving history (%s)", in.Username)

	for _, v := range results {
		stream.Send(&pb.RaceResult{
			Username:         v.Username,
			MotorcycleId:     int32(v.MotorcycleId),
			PositionInRace:   int32(v.Position),
			TotalMotorcycles: int32(v.TotalMotorcycles),
			TrackName:        v.TrackName,
			MotorcycleName:   v.MotorcycleName,
			MotorcycleLevel:  int32(v.MotorcycleLevel),
			Time:             timestamppb.New(v.Time),
		})
	}

	return nil
}

func (s *Server) StartMatchmaking(ctx context.Context, in *pb.RaceMotorcycle) (*emptypb.Empty, error) {
	track, left, err := s.db.StartMatchmaking(in.Username, &MotorcycleStats{
		Id:           int(in.MotorcycleId),
		Name:         in.MotorcycleName,
		Level:        int(in.Level),
		Engine:       int(in.Engine),
		Brakes:       int(in.Brakes),
		Aerodynamics: int(in.Aerodynamics),
		Agility:      int(in.Agility),
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Printf("Starting matchmaking (%s:%d)", in.Username, in.MotorcycleId)

	if left == 0 {
		results, err := s.db.CompleteRace(track)

		log.Printf("Race ended")

		if err != nil {
			log.Println(err)
			return nil, err
		}

		c := pb.NewOrchestratorClient(s.orchestrator)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		stream, err := c.NotifyEndRace(ctx)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		log.Printf("Sending results")

		for _, v := range results {
			err = stream.Send(&pb.RaceResult{
				Username:         v.Username,
				MotorcycleId:     int32(v.MotorcycleId),
				PositionInRace:   int32(v.Position),
				TotalMotorcycles: int32(v.TotalMotorcycles),
				TrackName:        v.TrackName,
				MotorcycleName:   v.MotorcycleName,
				MotorcycleLevel:  int32(v.MotorcycleLevel),
				Time:             timestamppb.New(v.Time),
			})

			if err != nil {
				log.Println(err)
				return nil, err
			}
		}

		stream.CloseSend()
		time.Sleep(time.Second)
	}

	return nil, err
}

func (s *Server) StillAlive(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, nil
}
