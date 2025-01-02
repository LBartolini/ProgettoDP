package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	pb "orchestrator/proto"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RaceResult struct {
	Username         string
	MotorcycleId     int
	MotorcycleName   string
	MotorcycleLevel  int
	Position         int
	TotalMotorcycles int
	TrackName        string
}

type LeaderboardInfo struct {
	Username string
	Points   int32
	Position int32
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
	Username              string
	MotorcycleId          int
	Name                  string
	Level                 int
	IsRacing              bool   // fetched from Racing Service
	TrackName             string // fetched from Racing Service
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

type Orchestrator struct {
	pb.UnimplementedOrchestratorServer
	// TODO separate in different kind of orchestrator based on service task
	balancer LoadBalancer // TODO ask directly to load balancer to do RPCs
}

func NewOrchestrator(balancer LoadBalancer) *Orchestrator {
	return &Orchestrator{balancer: balancer}
}

func getGrpcClientFromContext(ctx context.Context) (*grpc.ClientConn, error) {
	peerInfo, ok := peer.FromContext(ctx)
	if !ok {
		return nil, errors.New("unable to get peer information from context")
	}

	address := strings.Split(peerInfo.Addr.String(), ":")[0]
	return grpc.NewClient(fmt.Sprintf("%s:%s", address, os.Getenv("SERVICE_PORT")), grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func (o *Orchestrator) RegisterAuth(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	client, err := getGrpcClientFromContext(ctx)
	if err != nil {
		client.Close()
		log.Println(err)
		return nil, err
	}

	return nil, o.balancer.RegisterAuth(client)
}

func (o *Orchestrator) RegisterGarage(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	client, err := getGrpcClientFromContext(ctx)
	if err != nil {
		client.Close()
		log.Println(err)
		return nil, err
	}

	return nil, o.balancer.RegisterGarage(client)
}

func (o *Orchestrator) RegisterLeaderboard(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	client, err := getGrpcClientFromContext(ctx)
	if err != nil {
		client.Close()
		log.Println(err)
		return nil, err
	}

	return nil, o.balancer.RegisterLeaderboard(client)
}

func (o *Orchestrator) RegisterRacing(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	client, err := getGrpcClientFromContext(ctx)
	if err != nil {
		client.Close()
		log.Println(err)
		return nil, err
	}

	return nil, o.balancer.RegisterRacing(client)
}

func (o *Orchestrator) NotifyEndRace(stream pb.Orchestrator_NotifyEndRaceServer) error {
	log.Println("Race ended")
	for {
		race_result, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(nil)
		}
		if err != nil {
			log.Println(err)
			return err
		}

		// Garage: increase money
		garage_conn := o.balancer.GetGarage()
		if garage_conn == nil {
			return errors.New("unable to get connection to garage service")
		}

		garage_client := pb.NewGarageClient(garage_conn)
		ctxAlive, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		money_win, _ := strconv.Atoi(os.Getenv("MONEY_WIN"))
		money_last, _ := strconv.Atoi(os.Getenv("MONEY_LAST"))
		increase := &pb.MoneyIncrease{Username: race_result.Username,
			Money: int32(o.computeAfterRace(int(race_result.PositionInRace), int(race_result.TotalMotorcycles), money_win, money_last))}
		log.Printf("Race ended: increase money %s by %d", race_result.Username, increase.Money)
		_, err = garage_client.IncreaseUserMoney(ctxAlive, increase)
		if err != nil {
			log.Println(err)
			return err
		}

		// Leaderboard: increase points
		leaderboard_conn := o.balancer.GetLeaderboard()
		if leaderboard_conn == nil {
			return errors.New("unable to get connection to garage service")
		}

		leaderboard_client := pb.NewLeaderboardClient(leaderboard_conn)
		ctxAliveLeaderboard, cancel_leaderboard := context.WithTimeout(context.Background(), time.Second)
		defer cancel_leaderboard()

		points_win, _ := strconv.Atoi(os.Getenv("POINTS_WIN"))
		points_last, _ := strconv.Atoi(os.Getenv("POINTS_LAST"))
		points := &pb.PointIncrement{Username: race_result.Username,
			Points: int32(o.computeAfterRace(int(race_result.PositionInRace), int(race_result.TotalMotorcycles), points_win, points_last))}
		log.Printf("Race ended: increase points %s by %d", race_result.Username, points.Points)
		_, err = leaderboard_client.AddPoints(ctxAliveLeaderboard, points)
		if err != nil {
			log.Println(err)
			return err
		}

		cancel()
		cancel_leaderboard()
	}
}

//////////////////////

func (o *Orchestrator) Login(username string, password string) (res bool, e error) {
	conn := o.balancer.GetAuth()

	if conn == nil {
		return false, nil
	}

	c := pb.NewAuthenticationClient(conn)
	ctxAlive, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	result, err := c.Login(ctxAlive, &pb.PlayerCredentials{Username: username, Password: password})
	if err != nil {
		log.Println(err)
		return false, err
	}

	log.Printf("Login of %s result %t", username, result.Result)
	return result.Result, nil
}

func (o *Orchestrator) Register(username string, password string, email string, phone string) (res bool, e error) {
	// Register in Auth Service
	conn := o.balancer.GetAuth()

	if conn == nil {
		return false, nil
	}

	auth_client := pb.NewAuthenticationClient(conn)
	ctxAlive, cancel1 := context.WithTimeout(context.Background(), time.Second)
	defer cancel1()

	register_result, err := auth_client.Register(ctxAlive, &pb.PlayerDetails{Username: username, Password: password, Email: email, Phone: phone})
	if err != nil {
		log.Println(err)
		return false, err
	}
	log.Printf("Auth Register of %s result %t", username, register_result.Result)

	// Register in Leaderboard Service
	conn = o.balancer.GetLeaderboard()

	if conn == nil {
		return false, nil
	}

	leaderboard_client := pb.NewLeaderboardClient(conn)
	ctxAlive, cancel2 := context.WithTimeout(context.Background(), time.Second)
	defer cancel2()

	_, err = leaderboard_client.AddPoints(ctxAlive, &pb.PointIncrement{Username: username, Points: 0})
	if err != nil {
		log.Println(err)
		return false, err
	}

	// Register in Garage Service
	conn = o.balancer.GetGarage()

	if conn == nil {
		return false, nil
	}

	garage_client := pb.NewGarageClient(conn)
	ctxAlive, cancel3 := context.WithTimeout(context.Background(), time.Second)
	defer cancel3()

	start_money, err := strconv.Atoi(os.Getenv("START_MONEY"))
	if err != nil {
		log.Println(err)
		return false, err
	}

	_, err = garage_client.IncreaseUserMoney(ctxAlive, &pb.MoneyIncrease{Username: username, Money: int32(start_money)})
	if err != nil {
		log.Println(err)
		return false, err
	}

	// Register in Racing Service

	return register_result.Result, nil
}

func (o *Orchestrator) GetLeaderboardInfo(username string) (*LeaderboardInfo, error) {
	conn := o.balancer.GetLeaderboard()
	if conn == nil {
		return nil, errors.New("unable to get connection to leaderboard service")
	}

	leaderboard_client := pb.NewLeaderboardClient(conn)
	ctxAlive, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	pos, err := leaderboard_client.GetPlayer(ctxAlive, &pb.PlayerUsername{Username: username})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &LeaderboardInfo{Username: username, Points: pos.Points, Position: pos.Position}, nil
}

func (o *Orchestrator) GetFullLeaderboard() ([]*LeaderboardInfo, error) {
	conn := o.balancer.GetLeaderboard()
	if conn == nil {
		return nil, errors.New("unable to get connection to leaderboard service")
	}

	leaderboard_client := pb.NewLeaderboardClient(conn)
	ctxAlive, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := leaderboard_client.GetLeaderboard(ctxAlive, nil)
	if err != nil {
		return nil, err
	}

	var leaderboard []*LeaderboardInfo
	for {
		p, err := r.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println(err)
			return nil, err
		} else {
			pos := &LeaderboardInfo{Username: p.Username, Points: p.Points, Position: p.Position}
			leaderboard = append(leaderboard, pos)
		}
	}

	return leaderboard, nil
}

func (o *Orchestrator) GetRemainingMotorcycles(username string) ([]*Motorcycle, error) {
	conn := o.balancer.GetGarage()
	if conn == nil {
		return nil, errors.New("unable to get connection to garage service")
	}

	garage_client := pb.NewGarageClient(conn)
	ctxAlive, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := garage_client.GetRemainingMotorcycles(ctxAlive, &pb.PlayerUsername{Username: username})
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

func (o *Orchestrator) GetUserMotorcycles(username string) ([]*Ownership, error) {
	garage_conn := o.balancer.GetGarage()
	if garage_conn == nil {
		return nil, errors.New("unable to get connection to garage service")
	}
	garage_client := pb.NewGarageClient(garage_conn)

	racing_conn := o.balancer.GetRacing()
	if racing_conn == nil {
		return nil, errors.New("unable to get connection to racing service")
	}
	racing_client := pb.NewRacingClient(racing_conn)

	ctxAlive, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := garage_client.GetUserMotorcycles(ctxAlive, &pb.PlayerUsername{Username: username})
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
			ctxAliveRacing, cancel_racing := context.WithTimeout(context.Background(), time.Second)
			defer cancel_racing()

			status, err_racing := racing_client.CheckIsRacing(ctxAliveRacing, &pb.PlayerMotorcycle{Username: username, MotorcycleId: p.MotorcycleInfo.Id})
			if err_racing != nil {
				return nil, err_racing
			}

			pos := &Ownership{
				Username:              username,
				Level:                 int(p.Level),
				MotorcycleId:          int(p.MotorcycleInfo.Id),
				IsRacing:              status.IsRacing,
				TrackName:             status.TrackName,
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
			}

			motorcycles = append(motorcycles, pos)
		}
	}

	return motorcycles, nil
}

func (o *Orchestrator) GetUserMoney(username string) (int, error) {
	conn := o.balancer.GetGarage()
	if conn == nil {
		return 0, errors.New("unable to get connection to garage service")
	}

	garage_client := pb.NewGarageClient(conn)
	ctxAlive, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	money, err := garage_client.GetUserMoney(ctxAlive, &pb.PlayerUsername{Username: username})
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return int(money.Money), nil
}

func (o *Orchestrator) BuyMotorcycle(username string, MotorcycleId int) (bool, error) {
	conn := o.balancer.GetGarage()
	if conn == nil {
		return false, errors.New("unable to get connection to garage service")
	}

	garage_client := pb.NewGarageClient(conn)
	ctxAlive, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := garage_client.BuyMotorcycle(ctxAlive, &pb.PlayerMotorcycle{Username: username, MotorcycleId: int32(MotorcycleId)})
	if err != nil {
		log.Println(err)
		return false, err
	}

	return true, nil
}

func (o *Orchestrator) UpgradeMotorcycle(username string, MotorcycleId int) (bool, error) {
	conn := o.balancer.GetGarage()
	if conn == nil {
		return false, errors.New("unable to get connection to garage service")
	}

	garage_client := pb.NewGarageClient(conn)
	ctxAlive, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := garage_client.UpgradeMotorcycle(ctxAlive, &pb.PlayerMotorcycle{Username: username, MotorcycleId: int32(MotorcycleId)})
	if err != nil {
		log.Println(err)
		return false, err
	}

	return true, nil
}

func (o *Orchestrator) computeAfterRace(position int, total int, first int, last int) int {
	// line that goes from the point (1, first) to (total, last)

	m := (last - first) / (total - 1)
	return int(m*(position-1) + first)
}

func (o *Orchestrator) StartMatchmaking(username string, MotorcycleId int) error {
	conn := o.balancer.GetGarage()
	if conn == nil {
		return errors.New("unable to get connection to garage service")
	}

	garage_client := pb.NewGarageClient(conn)
	ctxAlive, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stats, err := garage_client.GetUserMotorcycleStats(ctxAlive, &pb.PlayerMotorcycle{Username: username, MotorcycleId: int32(MotorcycleId)})
	if err != nil {
		log.Println(err)
		return err
	}

	conn = o.balancer.GetRacing()
	if conn == nil {
		return errors.New("unable to get connection to racing service")
	}

	racing_client := pb.NewRacingClient(conn)
	ctxAlive, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = racing_client.StartMatchmaking(ctxAlive, &pb.RaceMotorcycle{Username: username,
		MotorcycleId:   int32(MotorcycleId),
		MotorcycleName: stats.MotorcycleInfo.Name,
		Level:          stats.Level,
		Engine:         stats.MotorcycleInfo.Engine,
		Brakes:         stats.MotorcycleInfo.Brakes,
		Agility:        stats.MotorcycleInfo.Agility,
		Aerodynamics:   stats.MotorcycleInfo.Aerodynamics})

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (o *Orchestrator) GetHistory(username string) ([]*RaceResult, error) {
	conn := o.balancer.GetRacing()
	if conn == nil {
		return nil, errors.New("unable to get connection to racing service")
	}

	racing_client := pb.NewRacingClient(conn)
	ctxAlive, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := racing_client.GetHistory(ctxAlive, &pb.PlayerUsername{Username: username})
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
				Username:         r.Username,
				MotorcycleId:     int(r.MotorcycleId),
				MotorcycleName:   r.MotorcycleName,
				MotorcycleLevel:  int(r.MotorcycleLevel),
				Position:         int(r.PositionInRace),
				TotalMotorcycles: int(r.TotalMotorcycles),
				TrackName:        r.TrackName}

			results = append(results, res)
		}
	}

	return results, nil
}
