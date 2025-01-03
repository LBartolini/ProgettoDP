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
		err = o.balancer.GarageIncreaseUserMoney(race_result.Username, increase)
		if err != nil {
			log.Println(err)
			return err
		}
		log.Printf("Race ended: increase money %s by %d", race_result.Username, increase)

		// Leaderboard: increase points
		points_win, _ := strconv.Atoi(os.Getenv("POINTS_WIN"))
		points_last, _ := strconv.Atoi(os.Getenv("POINTS_LAST"))
		increase = o.computeAfterRace(int(race_result.PositionInRace), int(race_result.TotalMotorcycles), points_win, points_last)
		err = o.balancer.LeaderboardAddPoints(race_result.Username, increase)
		if err != nil {
			log.Println(err)
			return err
		}
		log.Printf("Race ended: increase points %s by %d", race_result.Username, increase)
	}
}

//////////////////////

func (o *Orchestrator) Login(username string, password string) (res bool, e error) {
	result, err := o.balancer.AuthLogin(username, password)

	log.Printf("Login of %s result %t", username, result)
	return result, err
}

func (o *Orchestrator) Register(username string, password string, email string, phone string) (res bool, e error) {
	// Register in Auth Service
	register_result, err := o.balancer.AuthRegister(username, password, email, phone)
	if err != nil {
		log.Println(err)
		return false, err
	}
	log.Printf("Auth Register of %s result %t", username, register_result)

	// Register in Leaderboard Service
	err = o.balancer.LeaderboardAddPoints(username, 0)
	if err != nil {
		log.Println(err)
		return false, err
	}

	// Register in Garage Service
	start_money, _ := strconv.Atoi(os.Getenv("START_MONEY"))
	err = o.balancer.GarageIncreaseUserMoney(username, start_money)
	if err != nil {
		log.Println(err)
		return false, err
	}

	return register_result, nil
}

func (o *Orchestrator) GetLeaderboardInfo(username string) (*LeaderboardPosition, error) {
	return o.balancer.LeaderboardGetPlayer(username)
}

func (o *Orchestrator) GetFullLeaderboard() ([]*LeaderboardPosition, error) {
	return o.balancer.LeaderboardGetFullLeaderboard()
}

func (o *Orchestrator) GetRemainingMotorcycles(username string) ([]*Motorcycle, error) {
	return o.balancer.GarageGetRemainingMotorcycles(username)
}

func (o *Orchestrator) GetUserMotorcycles(username string) ([]*Ownership, error) {
	owned, err := o.balancer.GarageGetUserMotorcycles(username)
	if err != nil {
		return nil, err
	}

	for _, v := range owned {
		if status, err := o.balancer.RacingCheckIsRacing(username, v.Motorcycle.Id); err == nil {
			v.RacingStatus = status
		}
	}

	return owned, nil
}

func (o *Orchestrator) GetUserMoney(username string) (int, error) {
	return o.balancer.GarageGetUserMoney(username)
}

func (o *Orchestrator) BuyMotorcycle(username string, MotorcycleId int) error {
	return o.balancer.GarageBuyMotorcycle(username, MotorcycleId)
}

func (o *Orchestrator) UpgradeMotorcycle(username string, MotorcycleId int) error {
	return o.balancer.GarageUpgradeMotorcycle(username, MotorcycleId)
}

func (o *Orchestrator) computeAfterRace(position int, total int, first int, last int) int {
	// line that goes from the point (1, first) to (total, last)

	m := (last - first) / (total - 1)
	return int(m*(position-1) + first)
}

func (o *Orchestrator) StartMatchmaking(username string, MotorcycleId int) error {
	stats, err := o.balancer.GarageGetUserMotorcycleStats(username, MotorcycleId)
	if err != nil {
		log.Println(err)
		return err
	}

	return o.balancer.RacingStartMatchmaking(username, stats)
}

func (o *Orchestrator) GetHistory(username string) ([]*RaceResult, error) {
	return o.balancer.RacingGetHistory(username)
}
