package services

import (
	"context"
	"io"
	"log"
	pb "orchestrator/proto"
	"time"

	"google.golang.org/grpc"
)

type RacingStatus struct {
	Username     string
	MotorcycleId int
	Status       bool
	TrackName    string
}

type RaceResult struct {
	MotorcycleName   string
	MotorcycleLevel  int
	Position         int
	TotalMotorcycles int
	TrackName        string
	Time             time.Time
}

type Racing interface {
	RacingCheckIsRacing(username string, motorcycle_id int) (*RacingStatus, error)
	RacingStartMatchmaking(username string, stats *Ownership) error
	RacingGetHistory(username string) ([]*RaceResult, error)
}

type RacingService struct {
	conn *grpc.ClientConn
}

func NewRacingService(conn *grpc.ClientConn) *RacingService {
	return &RacingService{conn: conn}
}

func (s *RacingService) RacingStartMatchmaking(username string, stats *Ownership) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := pb.NewRacingClient(s.conn).StartMatchmaking(ctx, &pb.RaceMotorcycle{
		Username:       username,
		MotorcycleId:   int32(stats.Motorcycle.Id),
		MotorcycleName: stats.Motorcycle.Name,
		Level:          int32(stats.Level),
		Engine:         int32(stats.Motorcycle.Engine),
		Aerodynamics:   int32(stats.Motorcycle.Aerodynamics),
		Agility:        int32(stats.Motorcycle.Agility),
		Brakes:         int32(stats.Motorcycle.Brakes),
	})
	return err
}

func (s *RacingService) RacingCheckIsRacing(username string, motorcycle_id int) (*RacingStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	status, err := pb.NewRacingClient(s.conn).CheckIsRacing(ctx, &pb.PlayerMotorcycle{Username: username, MotorcycleId: int32(motorcycle_id)})
	if err != nil {
		return nil, err
	}

	return &RacingStatus{Username: username, MotorcycleId: motorcycle_id, Status: status.IsRacing, TrackName: status.TrackName}, nil
}

func (s *RacingService) RacingGetHistory(username string) ([]*RaceResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := pb.NewRacingClient(s.conn).GetHistory(ctx, &pb.PlayerUsername{Username: username})
	if err != nil {
		return nil, err
	}

	var results []*RaceResult
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println(err)
			return nil, err
		} else {
			res := &RaceResult{
				MotorcycleName:   r.MotorcycleName,
				MotorcycleLevel:  int(r.MotorcycleLevel),
				Position:         int(r.PositionInRace),
				TotalMotorcycles: int(r.TotalMotorcycles),
				TrackName:        r.TrackName,
				Time:             r.Time.AsTime(),
			}

			results = append(results, res)
		}
	}

	return results, nil
}
