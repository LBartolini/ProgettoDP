package internal

import (
	"context"
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

func (s *Server) GetRemainingMotorcycles(in *pb.PlayerUsername, stream pb.Garage_GetRemainingMotorcyclesServer) error {
	motorcycles, err := s.db.GetRemainingMotorcycles(in.Username)

	if len(motorcycles) == 0 || err != nil {
		return err
	}

	// TODO convert in RANGE
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

	if len(ownerships) == 0 || err != nil {
		return err
	}

	// TODO convert in RANGE
	for i := 0; i < len(ownerships); i++ {
		stream.Send(&pb.OwnershipInfo{
			Username:     ownerships[i].Username,
			MotorcycleId: int32(ownerships[i].MotorcycleId),
			Level:        int32(ownerships[i].Level),
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

func (s *Server) GetUserMotorcycleStats(ctx context.Context, in *pb.PlayerMotorcycle) (*pb.OwnershipInfo, error) {
	ownership, err := s.db.GetUserMotorcycleStats(in.Username, int(in.MotorcycleId))

	if err != nil {
		return nil, err
	}

	info := &pb.OwnershipInfo{
		Username:     ownership.Username,
		MotorcycleId: int32(ownership.MotorcycleId),
		Level:        int32(ownership.Level),
		MotorcycleInfo: &pb.MotorcycleInfo{
			Id:                    int32(ownership.MotorcycleId),
			Name:                  ownership.Name,
			PriceToBuy:            int32(ownership.PriceToBuy),
			PriceToUpgrade:        int32(ownership.PriceToUpgrade),
			MaxLevel:              int32(ownership.MaxLevel),
			Engine:                int32(ownership.Engine),
			EngineIncrement:       int32(ownership.EngineIncrement),
			Agility:               int32(ownership.Agility),
			AgilityIncrement:      int32(ownership.AgilityIncrement),
			Brakes:                int32(ownership.Brakes),
			BrakesIncrement:       int32(ownership.BrakesIncrement),
			Aerodynamics:          int32(ownership.Aerodynamics),
			AerodynamicsIncrement: int32(ownership.AerodynamicsIncrement),
		}}

	return info, nil
}

func (s *Server) GetUserMoney(ctx context.Context, in *pb.PlayerUsername) (*pb.UserMoney, error) {
	money, err := s.db.GetUserMoney(in.Username)

	return &pb.UserMoney{Money: int32(money)}, err
}

func (s *Server) IncreaseUserMoney(ctx context.Context, in *pb.MoneyIncrease) (*emptypb.Empty, error) {
	err := s.db.IncreaseUserMoney(in.Username, int(in.Money))

	return nil, err
}

func (s *Server) BuyMotorcycle(ctx context.Context, in *pb.PlayerMotorcycle) (*emptypb.Empty, error) {
	return nil, s.db.BuyMotorcycle(in.Username, int(in.MotorcycleId))
}

func (s *Server) UpgradeMotorcycle(ctx context.Context, in *pb.PlayerMotorcycle) (*emptypb.Empty, error) {
	return nil, s.db.UpgradeMotorcycle(in.Username, int(in.MotorcycleId))
}

func (s *Server) StillAlive(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	log.Printf("Still Alive")
	return nil, nil
}
