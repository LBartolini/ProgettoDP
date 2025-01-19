package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"orchestrator/internal/services"
	pb "orchestrator/proto"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Orchestrator struct {
	pb.UnimplementedOrchestratorServer
	balancer LoadBalancer
}

func NewOrchestrator(balancer LoadBalancer) *Orchestrator {
	return &Orchestrator{balancer: balancer}
}

func (o *Orchestrator) computeAfterRace(position int, total int, first int, last int) int {
	// Utility function that computes the line that goes from the point (1, first) to (total, last)

	m := (last - first) / (total - 1)
	return int(m*(position-1) + first)
}

func getGrpcClientFromContext(ctx context.Context) (*grpc.ClientConn, error) {
	// Utility function to obtain connection from context

	peerInfo, ok := peer.FromContext(ctx)
	if !ok {
		return nil, errors.New("unable to get peer information from context")
	}

	address := strings.Split(peerInfo.Addr.String(), ":")[0]
	return grpc.NewClient(fmt.Sprintf("%s:%s", address, os.Getenv("SERVICE_PORT")), grpc.WithTransportCredentials(insecure.NewCredentials()))
}

////////////////////// gRPC Server section

func (o *Orchestrator) RegisterAuth(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	client, err := getGrpcClientFromContext(ctx)
	if err != nil {
		client.Close()
		log.Println(err)
		return nil, err
	}

	o.balancer.RegisterAuth(services.NewAuthService(client))

	return nil, nil
}

func (o *Orchestrator) RegisterGarage(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	client, err := getGrpcClientFromContext(ctx)
	if err != nil {
		client.Close()
		log.Println(err)
		return nil, err
	}

	o.balancer.RegisterGarage(services.NewGarageService(client))

	return nil, nil
}

func (o *Orchestrator) RegisterLeaderboard(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	client, err := getGrpcClientFromContext(ctx)
	if err != nil {
		client.Close()
		log.Println(err)
		return nil, err
	}

	o.balancer.RegisterLeaderboard(services.NewLeaderboardService(client))

	return nil, nil
}

func (o *Orchestrator) RegisterRacing(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	client, err := getGrpcClientFromContext(ctx)
	if err != nil {
		client.Close()
		log.Println(err)
		return nil, err
	}

	o.balancer.RegisterRacing(services.NewRacingService(client))

	return nil, nil
}

func (o *Orchestrator) NotifyEndRace(stream pb.Orchestrator_NotifyEndRaceServer) error {
	// Used by Racing service to notify the end of a race by streaming results

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
		money_win, _ := strconv.Atoi(os.Getenv("MONEY_WIN"))
		money_last, _ := strconv.Atoi(os.Getenv("MONEY_LAST"))
		increase := o.computeAfterRace(int(race_result.PositionInRace), int(race_result.TotalMotorcycles), money_win, money_last)

		// Give money to user based on position in race
		garage := o.balancer.GetGarage()
		if garage == nil {
			return errors.New("unable to connect to Garage Service")
		}
		err = garage.IncreaseUserMoney(race_result.Username, increase)
		if err != nil {
			log.Println(err)
			return err
		}
		log.Printf("Race ended: increase money %s by %d", race_result.Username, increase)

		// Leaderboard: increase points
		leaderboard := o.balancer.GetLeaderboard()
		if leaderboard == nil {
			return errors.New("unable to connect to Leaderboard Service")
		}

		points_win, _ := strconv.Atoi(os.Getenv("POINTS_WIN"))
		points_last, _ := strconv.Atoi(os.Getenv("POINTS_LAST"))
		increase = o.computeAfterRace(int(race_result.PositionInRace), int(race_result.TotalMotorcycles), points_win, points_last)

		// Give points to user based on position in race
		err = leaderboard.AddPoints(race_result.Username, increase)
		if err != nil {
			log.Println(err)
			return err
		}
		log.Printf("Race ended: increase points %s by %d", race_result.Username, increase)
	}
}

////////////////////// Orchestrator section

func (o *Orchestrator) Login(username string, password string) (res bool, e error) {
	conn := o.balancer.GetAuth()

	if conn == nil {
		return false, errors.New("unable to connect to Auth Service")
	}

	// Proxy
	result, err := conn.Login(username, password)

	log.Printf("Login of %s result %t", username, result)
	return result, err
}

func (o *Orchestrator) Register(username string, password string, email string, phone string) (res bool, e error) {
	// Register in Auth Service
	auth := o.balancer.GetAuth()

	if auth == nil {
		return false, errors.New("unable to connect to Auth Service")
	}

	register_result, err := auth.Register(username, password, email, phone)
	if err != nil {
		log.Println(err)
		return false, err
	}
	log.Printf("Auth Register of %s result %t", username, register_result)

	// Register in Leaderboard Service (giving startin points)
	leaderboard := o.balancer.GetLeaderboard()

	if leaderboard == nil {
		return false, errors.New("unable to connect to Leaderboard Service")
	}

	err = leaderboard.AddPoints(username, 0)
	if err != nil {
		log.Println(err)
		return false, err
	}

	// Register in Garage Service (giving starting money)
	garage := o.balancer.GetGarage()

	if garage == nil {
		return false, errors.New("unable to connect to Garage Service")
	}

	start_money, _ := strconv.Atoi(os.Getenv("START_MONEY"))
	err = garage.IncreaseUserMoney(username, start_money)
	if err != nil {
		log.Println(err)
		return false, err
	}

	return register_result, nil
}

func (o *Orchestrator) GetLeaderboardInfo(username string) (*services.LeaderboardPosition, error) {
	conn := o.balancer.GetLeaderboard()

	if conn == nil {
		return nil, errors.New("unable to connect to Leaderboard Service")
	}

	// Proxy
	return conn.GetPlayer(username)
}

func (o *Orchestrator) GetFullLeaderboard() ([]*services.LeaderboardPosition, error) {
	conn := o.balancer.GetLeaderboard()

	if conn == nil {
		return nil, errors.New("unable to connect to Leaderboard Service")
	}

	// Proxy
	return conn.GetFullLeaderboard()
}

func (o *Orchestrator) GetRemainingMotorcycles(username string) ([]*services.Motorcycle, error) {
	garage_conn := o.balancer.GetGarage()

	if garage_conn == nil {
		return nil, errors.New("unable to connect to Garage Service")
	}

	// Proxy
	return garage_conn.GetRemainingMotorcycles(username)
}

func (o *Orchestrator) GetUserMotorcycles(username string) ([]*services.Ownership, error) {
	garage_conn := o.balancer.GetGarage()

	if garage_conn == nil {
		return nil, errors.New("unable to connect to Garage Service")
	}

	owned, err := garage_conn.GetUserMotorcycles(username)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	racing_conn := o.balancer.GetRacing()

	if racing_conn == nil {
		return nil, errors.New("unable to connect to Racing Service")
	}

	for _, v := range owned {
		// Add information regarding racing to each motorcycle
		if status, err := racing_conn.CheckIsRacing(username, v.Motorcycle.Id); err == nil {
			v.RacingStatus = status
		}
	}

	return owned, nil
}

func (o *Orchestrator) GetUserMoney(username string) (int, error) {
	garage_conn := o.balancer.GetGarage()

	if garage_conn == nil {
		return 0, errors.New("unable to connect to Garage Service")
	}

	// Proxy
	return garage_conn.GetUserMoney(username)
}

func (o *Orchestrator) BuyMotorcycle(username string, MotorcycleId int) error {
	garage_conn := o.balancer.GetGarage()

	if garage_conn == nil {
		return errors.New("unable to connect to Garage Service")
	}

	// Proxy
	return garage_conn.BuyMotorcycle(username, MotorcycleId)
}

func (o *Orchestrator) UpgradeMotorcycle(username string, MotorcycleId int) error {
	garage_conn := o.balancer.GetGarage()

	if garage_conn == nil {
		return errors.New("unable to connect to Garage Service")
	}

	// Proxy
	return garage_conn.UpgradeMotorcycle(username, MotorcycleId)
}

func (o *Orchestrator) StartMatchmaking(username string, MotorcycleId int) error {
	garage_conn := o.balancer.GetGarage()

	if garage_conn == nil {
		return errors.New("unable to connect to Garage Service")
	}

	// Get stats from garage
	stats, err := garage_conn.GetUserMotorcycleStats(username, MotorcycleId)
	if err != nil {
		log.Println(err)
		return err
	}

	racing_conn := o.balancer.GetRacing()

	if racing_conn == nil {
		return errors.New("unable to connect to Racing Service")
	}

	// Start matchmaking sending stats of motorcycle
	return racing_conn.StartMatchmaking(username, stats)
}

func (o *Orchestrator) GetHistory(username string) ([]*services.RaceResult, error) {
	conn := o.balancer.GetRacing()

	if conn == nil {
		return nil, errors.New("unable to connect to Racing Service")
	}

	// Proxy
	return conn.GetHistory(username)
}
