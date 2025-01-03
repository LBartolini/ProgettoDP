package internal

import (
	"context"
	"errors"
	"io"
	"log"
	pb "orchestrator/proto"
	"sync"
	"time"

	"math/rand/v2"

	"google.golang.org/grpc"
)

type LeaderboardPosition struct {
	Username string
	Points   int
	Position int
}

type RacingStatus struct {
	Username     string
	MotorcycleId int
	Status       bool
	TrackName    string
}

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

type RaceResult struct {
	MotorcycleName   string
	MotorcycleLevel  int
	Position         int
	TotalMotorcycles int
	TrackName        string
	Time             time.Time
}

type LoadBalancer interface {
	RegisterAuth(conn *grpc.ClientConn) error
	GetAuth() *grpc.ClientConn

	RegisterLeaderboard(conn *grpc.ClientConn) error
	GetLeaderboard() *grpc.ClientConn

	RegisterGarage(conn *grpc.ClientConn) error
	GetGarage() *grpc.ClientConn

	RegisterRacing(conn *grpc.ClientConn) error
	GetRacing() *grpc.ClientConn

	testConnection(*grpc.ClientConn) error

	AuthLogin(username string, password string) (bool, error)
	AuthRegister(username, password, email, phone string) (bool, error)

	GarageGetUserMoney(username string) (int, error)
	GarageIncreaseUserMoney(username string, money int) error
	GarageGetRemainingMotorcycles(username string) ([]*Motorcycle, error)
	GarageGetUserMotorcycles(username string) ([]*Ownership, error)
	GarageBuyMotorcycle(username string, motorcycle_id int) error
	GarageUpgradeMotorcycle(username string, motorcycle_id int) error
	GarageGetUserMotorcycleStats(username string, motorcycle_id int) (*Ownership, error)

	LeaderboardAddPoints(username string, points int) error
	LeaderboardGetPlayer(username string) (*LeaderboardPosition, error)
	LeaderboardGetFullLeaderboard() ([]*LeaderboardPosition, error)

	RacingCheckIsRacing(username string, motorcycle_id int) (*RacingStatus, error)
	RacingStartMatchmaking(username string, stats *Ownership) error
	RacingGetHistory(username string) ([]*RaceResult, error)
}

type RandomLoadBalancer struct {
	mu       sync.Mutex
	services map[string][]*grpc.ClientConn
}

func NewRandomLoadBalancer() *RandomLoadBalancer {
	return &RandomLoadBalancer{services: make(map[string][]*grpc.ClientConn)}
}

func (lb *RandomLoadBalancer) registerService(name string, conn *grpc.ClientConn) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if _, exists := lb.services[name]; !exists {
		lb.services[name] = []*grpc.ClientConn{}
	}
	lb.services[name] = append(lb.services[name], conn)
	return nil
}

func removeAtIndex(slice []*grpc.ClientConn, index int) []*grpc.ClientConn {
	if index < 0 || index >= len(slice) {
		return slice
	}
	return append(slice[:index], slice[index+1:]...)
}

func (lb *RandomLoadBalancer) getService(name string) *grpc.ClientConn {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	var conn *grpc.ClientConn

	for i := len(lb.services[name]); i > 0; i-- {
		index := rand.IntN(len(lb.services[name]))
		temp := lb.services[name][index]

		if err := lb.testConnection(temp); err == nil {
			conn = temp
			log.Printf("Service %s found", name)
			break
		} else {
			log.Printf("Balancer removing service %s not alive", name)
			lb.services[name] = removeAtIndex(lb.services[name], index)
			temp.Close()
		}
	}

	return conn
}

func (lb *RandomLoadBalancer) RegisterAuth(conn *grpc.ClientConn) error {
	return lb.registerService("auth", conn)
}

func (lb *RandomLoadBalancer) GetAuth() *grpc.ClientConn {
	return lb.getService("auth")
}

func (lb *RandomLoadBalancer) RegisterLeaderboard(conn *grpc.ClientConn) error {
	return lb.registerService("leaderboard", conn)
}

func (lb *RandomLoadBalancer) GetLeaderboard() *grpc.ClientConn {
	return lb.getService("leaderboard")
}

func (lb *RandomLoadBalancer) RegisterRacing(conn *grpc.ClientConn) error {
	return lb.registerService("racing", conn)
}

func (lb *RandomLoadBalancer) GetRacing() *grpc.ClientConn {
	return lb.getService("racing")
}

func (lb *RandomLoadBalancer) RegisterGarage(conn *grpc.ClientConn) error {
	return lb.registerService("garage", conn)
}

func (lb *RandomLoadBalancer) GetGarage() *grpc.ClientConn {
	return lb.getService("garage")
}

func (lb *RandomLoadBalancer) testConnection(conn *grpc.ClientConn) error {
	c := pb.NewStillAliveClient(conn)
	ctxAlive, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := c.StillAlive(ctxAlive, nil)
	return err
}

func (lb *RandomLoadBalancer) AuthLogin(username string, password string) (bool, error) {
	conn := lb.GetAuth()

	if conn == nil {
		return false, errors.New("unable to connect to Auth Service")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := pb.NewAuthenticationClient(conn).Login(ctx, &pb.PlayerCredentials{Username: username, Password: password})
	if err != nil || res == nil {
		return false, err
	}

	return res.Result, nil
}

func (lb *RandomLoadBalancer) AuthRegister(username, password, email, phone string) (bool, error) {
	conn := lb.GetAuth()

	if conn == nil {
		return false, errors.New("unable to connect to Auth Service")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := pb.NewAuthenticationClient(conn).Register(ctx, &pb.PlayerDetails{Username: username, Password: password, Email: email, Phone: phone})
	if err != nil || res == nil {
		return false, err
	}

	return res.Result, nil
}

func (lb *RandomLoadBalancer) GarageGetRemainingMotorcycles(username string) ([]*Motorcycle, error) {
	conn := lb.GetGarage()

	if conn == nil {
		return nil, errors.New("unable to connect to Garage Service")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := pb.NewGarageClient(conn).GetRemainingMotorcycles(ctx, &pb.PlayerUsername{Username: username})
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

func (lb *RandomLoadBalancer) GarageGetUserMotorcycles(username string) ([]*Ownership, error) {
	conn := lb.GetGarage()

	if conn == nil {
		return nil, errors.New("unable to connect to Garage Service")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := pb.NewGarageClient(conn).GetUserMotorcycles(ctx, &pb.PlayerUsername{Username: username})
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

func (lb *RandomLoadBalancer) GarageGetUserMotorcycleStats(username string, motorcycle_id int) (*Ownership, error) {
	conn := lb.GetGarage()

	if conn == nil {
		return nil, errors.New("unable to connect to Garage Service")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	p, err := pb.NewGarageClient(conn).GetUserMotorcycleStats(ctx, &pb.PlayerMotorcycle{Username: username, MotorcycleId: int32(motorcycle_id)})
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

func (lb *RandomLoadBalancer) GarageGetUserMoney(username string) (int, error) {
	conn := lb.GetGarage()

	if conn == nil {
		return 0, errors.New("unable to connect to Garage Service")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := pb.NewGarageClient(conn).GetUserMoney(ctx, &pb.PlayerUsername{Username: username})
	if err != nil {
		return 0, err
	}

	return int(res.Money), nil
}

func (lb *RandomLoadBalancer) GarageIncreaseUserMoney(username string, money int) error {
	conn := lb.GetGarage()

	if conn == nil {
		return errors.New("unable to connect to Garage Service")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := pb.NewGarageClient(conn).IncreaseUserMoney(ctx, &pb.MoneyIncrease{Username: username, Money: int32(money)})
	return err
}

func (lb *RandomLoadBalancer) GarageBuyMotorcycle(username string, motorcycle_id int) error {
	conn := lb.GetGarage()

	if conn == nil {
		return errors.New("unable to connect to Garage Service")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := pb.NewGarageClient(conn).BuyMotorcycle(ctx, &pb.PlayerMotorcycle{Username: username, MotorcycleId: int32(motorcycle_id)})
	return err
}

func (lb *RandomLoadBalancer) GarageUpgradeMotorcycle(username string, motorcycle_id int) error {
	conn := lb.GetGarage()

	if conn == nil {
		return errors.New("unable to connect to Garage Service")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := pb.NewGarageClient(conn).UpgradeMotorcycle(ctx, &pb.PlayerMotorcycle{Username: username, MotorcycleId: int32(motorcycle_id)})
	return err
}

func (lb *RandomLoadBalancer) LeaderboardGetFullLeaderboard() ([]*LeaderboardPosition, error) {
	conn := lb.GetLeaderboard()

	if conn == nil {
		return nil, errors.New("unable to connect to Leaderboard Service")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := pb.NewLeaderboardClient(conn).GetFullLeaderboard(ctx, nil)
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

func (lb *RandomLoadBalancer) LeaderboardGetPlayer(username string) (*LeaderboardPosition, error) {
	conn := lb.GetLeaderboard()

	if conn == nil {
		return nil, errors.New("unable to connect to Leaderboard Service")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	pos, err := pb.NewLeaderboardClient(conn).GetPlayer(ctx, &pb.PlayerUsername{Username: username})
	if err != nil {
		return nil, err
	}
	return &LeaderboardPosition{Username: pos.Username, Points: int(pos.Points), Position: int(pos.Position)}, nil
}

func (lb *RandomLoadBalancer) LeaderboardAddPoints(username string, points int) error {
	conn := lb.GetLeaderboard()

	if conn == nil {
		return errors.New("unable to connect to Leaderboard Service")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := pb.NewLeaderboardClient(conn).AddPoints(ctx, &pb.PointIncrement{Username: username, Points: int32(points)})
	return err
}

func (lb *RandomLoadBalancer) RacingStartMatchmaking(username string, stats *Ownership) error {
	conn := lb.GetRacing()

	if conn == nil {
		return errors.New("unable to connect to Racing Service")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := pb.NewRacingClient(conn).StartMatchmaking(ctx, &pb.RaceMotorcycle{
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

func (lb *RandomLoadBalancer) RacingCheckIsRacing(username string, motorcycle_id int) (*RacingStatus, error) {
	conn := lb.GetRacing()

	if conn == nil {
		return nil, errors.New("unable to connect to Racing Service")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	status, err := pb.NewRacingClient(conn).CheckIsRacing(ctx, &pb.PlayerMotorcycle{Username: username, MotorcycleId: int32(motorcycle_id)})
	if err != nil {
		return nil, err
	}

	return &RacingStatus{Username: username, MotorcycleId: motorcycle_id, Status: status.IsRacing, TrackName: status.TrackName}, nil
}

func (lb *RandomLoadBalancer) RacingGetHistory(username string) ([]*RaceResult, error) {
	conn := lb.GetRacing()

	if conn == nil {
		return nil, errors.New("unable to connect to Racing Service")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := pb.NewRacingClient(conn).GetHistory(ctx, &pb.PlayerUsername{Username: username})
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
