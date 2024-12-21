package internal

import (
	"context"
	"errors"
	"log"

	pb "garage/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	pb.UnimplementedGarageServer
	pb.UnimplementedStillAliveServer
	db GarageDB
}

func NewServer(conn GarageDB) *Server {
	return &Server{db: conn}
}

func (s *Server) GetAllMotorcycles(_ *emptypb.Empty, stream pb.Garage_GetAllMotorcyclesServer) error {
	motorcycles := s.db.GetAllMotorcycles()

	if len(motorcycles) == 0 || motorcycles == nil {
		return errors.New("no motorcycles")
	}

	for i := 0; i < len(motorcycles); i++ {
		stream.Send(&pb.MotorcycleInfo{
			Id:                    int32(motorcycles[i].Id),
			Name:                  motorcycles[i].Name,
			PriceToBuy:            int32(motorcycles[i].PriceToBuy),
			PriceToUpgrade:        int32(motorcycles[i].PriceToUpgrade),
			MaxLevel:              int32(motorcycles[i].MaxLevel),
			Engine:                int32(motorcycles[i].Engine),
			EngineIncrement:       int32(motorcycles[i].EngineIncrement),
			Agility:               int32(motorcycles[i].Agility),
			AgilityIncrement:      int32(motorcycles[i].AgilityIncrement),
			Brakes:                int32(motorcycles[i].Brakes),
			BrakesIncrement:       int32(motorcycles[i].BrakesIncrement),
			Aerodynamics:          int32(motorcycles[i].Aerodynamics),
			AerodynamicsIncrement: int32(motorcycles[i].AerodynamicsIncrement),
		})
	}

	return nil
}

func (s *Server) GetUserMotorcycles(in *pb.PlayerUsername, stream pb.Garage_GetUserMotorcyclesServer) error {
	ownerships, err := s.db.GetUserMotorcycles(in.Username)

	if len(ownerships) == 0 || ownerships == nil || err != nil {
		return errors.New("no motorcycles in user garage")
	}

	for i := 0; i < len(ownerships); i++ {
		stream.Send(&pb.OwnershipInfo{
			Username:     ownerships[i].Username,
			MotorcycleId: int32(ownerships[i].MotorcycleId),
			Level:        int32(ownerships[i].Level),
			IsRacing:     ownerships[i].IsRacing,
			MotorcycleInfo: &pb.MotorcycleInfo{
				Id:                    int32(ownerships[i].MotorcycleId),
				Name:                  ownerships[i].Name,
				PriceToBuy:            int32(ownerships[i].PriceToBuy),
				PriceToUpgrade:        int32(ownerships[i].PriceToUpgrade),
				MaxLevel:              int32(ownerships[i].MaxLevel),
				Engine:                int32(ownerships[i].Engine),
				EngineIncrement:       int32(ownerships[i].EngineIncrement),
				Agility:               int32(ownerships[i].Agility),
				AgilityIncrement:      int32(ownerships[i].AgilityIncrement),
				Brakes:                int32(ownerships[i].Brakes),
				BrakesIncrement:       int32(ownerships[i].BrakesIncrement),
				Aerodynamics:          int32(ownerships[i].Aerodynamics),
				AerodynamicsIncrement: int32(ownerships[i].AerodynamicsIncrement),
			},
		})
	}

	return nil
}

func (s *Server) GetUserMoney(ctx context.Context, in *pb.PlayerUsername) (*pb.UserMoney, error) {
	money, err := s.db.GetUserMoney(in.Username)

	return &pb.UserMoney{Money: int32(money)}, err
}

func (s *Server) BuyMotorcycle(ctx context.Context, in *pb.PlayerMotorcycle) (*emptypb.Empty, error) {
	return nil, s.db.BuyMotorcycle(in.Username, int(in.MotorcycleId))
}

func (s *Server) UpgradeMotorcycle(ctx context.Context, in *pb.PlayerMotorcycle) (*emptypb.Empty, error) {
	return nil, s.db.UpgradeMotorcycle(in.Username, int(in.MotorcycleId))
}

func (s *Server) StartRace(ctx context.Context, in *pb.PlayerMotorcycle) (*emptypb.Empty, error) {
	return nil, s.db.StartRace(in.Username, int(in.MotorcycleId))
}

func (s *Server) EndRace(ctx context.Context, in *pb.PlayerMotorcycle) (*emptypb.Empty, error) {
	return nil, s.db.EndRace(in.Username, int(in.MotorcycleId))
}

func (s *Server) StillAlive(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	log.Printf("Still Alive")
	return nil, nil
}