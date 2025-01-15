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
		log.Println(err)
		return err
	}

	log.Printf("Retrieving motorcycles not owned (%s)", in.Username)

	for _, v := range motorcycles {
		stream.Send(&pb.MotorcycleInfo{
			Id:                    int32(v.Id),
			Name:                  v.Name,
			PriceToBuy:            int32(v.PriceToBuy),
			PriceToUpgrade:        int32(v.PriceToUpgrade),
			MaxLevel:              int32(v.MaxLevel),
			Engine:                int32(v.Engine),
			EngineIncrement:       int32(v.EngineIncrement),
			Agility:               int32(v.Agility),
			AgilityIncrement:      int32(v.AgilityIncrement),
			Brakes:                int32(v.Brakes),
			BrakesIncrement:       int32(v.BrakesIncrement),
			Aerodynamics:          int32(v.Aerodynamics),
			AerodynamicsIncrement: int32(v.AerodynamicsIncrement),
		})
	}

	return nil
}

func (s *Server) GetUserMotorcycles(in *pb.PlayerUsername, stream pb.Garage_GetUserMotorcyclesServer) error {
	ownerships, err := s.db.GetUserMotorcycles(in.Username)

	if len(ownerships) == 0 || err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Retrieving motorcycles owned (%s)", in.Username)

	for _, v := range ownerships {
		stream.Send(&pb.OwnershipInfo{
			Username:     v.Username,
			MotorcycleId: int32(v.MotorcycleId),
			Level:        int32(v.Level),
			MotorcycleInfo: &pb.MotorcycleInfo{
				Id:                    int32(v.MotorcycleId),
				Name:                  v.Name,
				PriceToBuy:            int32(v.PriceToBuy),
				PriceToUpgrade:        int32(v.PriceToUpgrade),
				MaxLevel:              int32(v.MaxLevel),
				Engine:                int32(v.Engine),
				EngineIncrement:       int32(v.EngineIncrement),
				Agility:               int32(v.Agility),
				AgilityIncrement:      int32(v.AgilityIncrement),
				Brakes:                int32(v.Brakes),
				BrakesIncrement:       int32(v.BrakesIncrement),
				Aerodynamics:          int32(v.Aerodynamics),
				AerodynamicsIncrement: int32(v.AerodynamicsIncrement),
			},
		})
	}

	return nil
}

func (s *Server) GetUserMotorcycleStats(ctx context.Context, in *pb.PlayerMotorcycle) (*pb.OwnershipInfo, error) {
	ownership, err := s.db.GetUserMotorcycleStats(in.Username, int(in.MotorcycleId))

	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Printf("Motorcyce stats (%s:%d)", in.Username, in.MotorcycleId)

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

	log.Printf("Retrieving money (%s) with value %d", in.Username, money)

	return &pb.UserMoney{Money: int32(money)}, err
}

func (s *Server) IncreaseUserMoney(ctx context.Context, in *pb.MoneyIncrease) (*emptypb.Empty, error) {
	err := s.db.IncreaseUserMoney(in.Username, int(in.Money))

	log.Printf("Increasing money (%s) with value %d", in.Username, in.Money)

	return nil, err
}

func (s *Server) BuyMotorcycle(ctx context.Context, in *pb.PlayerMotorcycle) (*emptypb.Empty, error) {
	log.Printf("Buying motorcycle (%s:%d)", in.Username, in.MotorcycleId)

	return nil, s.db.BuyMotorcycle(in.Username, int(in.MotorcycleId))
}

func (s *Server) UpgradeMotorcycle(ctx context.Context, in *pb.PlayerMotorcycle) (*emptypb.Empty, error) {
	log.Printf("Upgrading motorcycle (%s:%d)", in.Username, in.MotorcycleId)

	return nil, s.db.UpgradeMotorcycle(in.Username, int(in.MotorcycleId))
}

func (s *Server) StillAlive(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, nil
}
