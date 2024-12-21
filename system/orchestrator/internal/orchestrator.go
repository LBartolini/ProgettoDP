package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	pb "orchestrator/proto"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/emptypb"
)

type LeaderboardInfo struct {
	Username string
	Points   int32
	Position int32
}

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
	ctxAlive, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

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
	ctxAlive, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = leaderboard_client.AddPoints(ctxAlive, &pb.PointIncrement{Username: username, Points: 0})
	if err != nil {
		log.Println(err)
		return false, err
	}

	// TODO: Register in Garage

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
