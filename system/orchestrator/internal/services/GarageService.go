package services

import (
	"context"
	"io"
	"log"
	"time"

	pb "orchestrator/proto"

	"google.golang.org/grpc"
)

type Motorcycle struct {
	Id                    int
	Name                  string
	PriceToBuy            int
	PriceToUpgrade        int
	MaxLevel              int
	Engine                int
	EngineIncrement       int
	Agility               int
	AgilityIncrement      int
	Brakes                int
	BrakesIncrement       int
	Aerodynamics          int
	AerodynamicsIncrement int
}

type Ownership struct {
	Level        int
	Motorcycle   *Motorcycle
	RacingStatus *RacingStatus
}

type Garage interface {
	StillAlive
	GetUserMoney(username string) (int, error)
	IncreaseUserMoney(username string, money int) error
	GetRemainingMotorcycles(username string) ([]*Motorcycle, error)
	GetUserMotorcycles(username string) ([]*Ownership, error)
	BuyMotorcycle(username string, motorcycle_id int) error
	UpgradeMotorcycle(username string, motorcycle_id int) error
	GetUserMotorcycleStats(username string, motorcycle_id int) (*Ownership, error)
}

// gRPC implementation of Garage interface
type GarageService struct {
	conn *grpc.ClientConn
}

func NewGarageService(conn *grpc.ClientConn) *GarageService {
	return &GarageService{conn: conn}
}

func (s *GarageService) StillAlive() bool {
	return StillAliveHandle(s.conn)
}

func (s *GarageService) Close() {
	s.conn.Close()
}

func (s *GarageService) GetRemainingMotorcycles(username string) ([]*Motorcycle, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := pb.NewGarageClient(s.conn).GetRemainingMotorcycles(ctx, &pb.PlayerUsername{Username: username})
	if err != nil {
		return nil, err
	}

	var motorcycles []*Motorcycle
	for {
		p, err := r.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println(err)
			return nil, err
		} else {
			pos := &Motorcycle{
				Id:                    int(p.Id),
				Name:                  p.Name,
				PriceToBuy:            int(p.PriceToBuy),
				PriceToUpgrade:        int(p.PriceToUpgrade),
				MaxLevel:              int(p.MaxLevel),
				Engine:                int(p.Engine),
				EngineIncrement:       int(p.EngineIncrement),
				Agility:               int(p.Agility),
				AgilityIncrement:      int(p.AgilityIncrement),
				Brakes:                int(p.Brakes),
				BrakesIncrement:       int(p.BrakesIncrement),
				Aerodynamics:          int(p.Aerodynamics),
				AerodynamicsIncrement: int(p.AerodynamicsIncrement),
			}
			motorcycles = append(motorcycles, pos)
		}
	}

	return motorcycles, nil
}

func (s *GarageService) GetUserMotorcycles(username string) ([]*Ownership, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := pb.NewGarageClient(s.conn).GetUserMotorcycles(ctx, &pb.PlayerUsername{Username: username})
	if err != nil {
		return nil, err
	}

	var motorcycles []*Ownership
	for {
		p, err := r.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println(err)
			return nil, err
		} else {
			m := &Ownership{
				Level: int(p.Level),
				Motorcycle: &Motorcycle{
					Id:                    int(p.MotorcycleInfo.Id),
					Name:                  p.MotorcycleInfo.Name,
					PriceToBuy:            int(p.MotorcycleInfo.PriceToBuy),
					PriceToUpgrade:        int(p.MotorcycleInfo.PriceToUpgrade),
					MaxLevel:              int(p.MotorcycleInfo.MaxLevel),
					Engine:                int(p.MotorcycleInfo.Engine),
					EngineIncrement:       int(p.MotorcycleInfo.EngineIncrement),
					Agility:               int(p.MotorcycleInfo.Agility),
					AgilityIncrement:      int(p.MotorcycleInfo.AgilityIncrement),
					Brakes:                int(p.MotorcycleInfo.Brakes),
					BrakesIncrement:       int(p.MotorcycleInfo.BrakesIncrement),
					Aerodynamics:          int(p.MotorcycleInfo.Aerodynamics),
					AerodynamicsIncrement: int(p.MotorcycleInfo.AerodynamicsIncrement),
				},
				RacingStatus: nil, // defined by racing service
			}

			motorcycles = append(motorcycles, m)
		}
	}

	return motorcycles, nil
}

func (s *GarageService) GetUserMotorcycleStats(username string, motorcycle_id int) (*Ownership, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	p, err := pb.NewGarageClient(s.conn).GetUserMotorcycleStats(ctx, &pb.PlayerMotorcycle{Username: username, MotorcycleId: int32(motorcycle_id)})
	if err != nil {
		return nil, err
	}

	stats := &Ownership{
		Level: int(p.Level),
		Motorcycle: &Motorcycle{
			Id:                    int(p.MotorcycleId),
			Name:                  p.MotorcycleInfo.Name,
			PriceToBuy:            int(p.MotorcycleInfo.PriceToBuy),
			PriceToUpgrade:        int(p.MotorcycleInfo.PriceToUpgrade),
			MaxLevel:              int(p.MotorcycleInfo.MaxLevel),
			Engine:                int(p.MotorcycleInfo.Engine),
			EngineIncrement:       int(p.MotorcycleInfo.EngineIncrement),
			Agility:               int(p.MotorcycleInfo.Agility),
			AgilityIncrement:      int(p.MotorcycleInfo.AgilityIncrement),
			Brakes:                int(p.MotorcycleInfo.Brakes),
			BrakesIncrement:       int(p.MotorcycleInfo.BrakesIncrement),
			Aerodynamics:          int(p.MotorcycleInfo.Aerodynamics),
			AerodynamicsIncrement: int(p.MotorcycleInfo.AerodynamicsIncrement),
		},
	}

	return stats, nil
}

func (s *GarageService) GetUserMoney(username string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := pb.NewGarageClient(s.conn).GetUserMoney(ctx, &pb.PlayerUsername{Username: username})
	if err != nil {
		return 0, err
	}

	return int(res.Money), nil
}

func (s *GarageService) IncreaseUserMoney(username string, money int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := pb.NewGarageClient(s.conn).IncreaseUserMoney(ctx, &pb.MoneyIncrease{Username: username, Money: int32(money)})
	return err
}

func (s *GarageService) BuyMotorcycle(username string, motorcycle_id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := pb.NewGarageClient(s.conn).BuyMotorcycle(ctx, &pb.PlayerMotorcycle{Username: username, MotorcycleId: int32(motorcycle_id)})
	return err
}

func (s *GarageService) UpgradeMotorcycle(username string, motorcycle_id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := pb.NewGarageClient(s.conn).UpgradeMotorcycle(ctx, &pb.PlayerMotorcycle{Username: username, MotorcycleId: int32(motorcycle_id)})
	return err
}
