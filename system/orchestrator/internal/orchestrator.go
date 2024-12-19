package internal

import (
	"context"
	"errors"
	"fmt"
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

type MyOrchestrator struct {
	pb.UnimplementedOrchestratorServer
	balancer LoadBalancer
}

func NewMyOrchestrator(balancer LoadBalancer) *MyOrchestrator {
	return &MyOrchestrator{balancer: balancer}
}

func getGrpcClientFromContext(ctx context.Context) (*grpc.ClientConn, error) {
	peerInfo, ok := peer.FromContext(ctx)
	if !ok {
		return nil, errors.New("unable to get peer information from context")
	}

	address := strings.Split(peerInfo.Addr.String(), ":")[0]
	return grpc.NewClient(fmt.Sprintf("%s:%s", address, os.Getenv("SERVICE_PORT")), grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func (o *MyOrchestrator) RegisterAuth(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	client, err := getGrpcClientFromContext(ctx)
	if err != nil {
		client.Close()
		return nil, err
	}

	return nil, o.balancer.RegisterAuth(client)
}

func (o *MyOrchestrator) RegisterGarage(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	client, err := getGrpcClientFromContext(ctx)
	if err != nil {
		client.Close()
		return nil, err
	}

	return nil, o.balancer.RegisterGarage(client)
}

func (o *MyOrchestrator) RegisterLeaderboard(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	client, err := getGrpcClientFromContext(ctx)
	if err != nil {
		client.Close()
		return nil, err
	}

	return nil, o.balancer.RegisterLeaderboard(client)
}

func (o *MyOrchestrator) RegisterRacing(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	client, err := getGrpcClientFromContext(ctx)
	if err != nil {
		client.Close()
		return nil, err
	}

	return nil, o.balancer.RegisterRacing(client)
}

func (o *MyOrchestrator) Login(username string, password string) (bool, error) {
	conn := o.balancer.GetAuth()

	if conn == nil {
		return false, nil
	}

	c := pb.NewAuthenticationClient(conn)
	ctxAlive, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	result, err := c.Login(ctxAlive, &pb.PlayerCredentials{Username: username, Password: password})
	if err != nil {
		return false, err
	}

	log.Printf("Login of %s result %t", username, result.Result)
	return result.Result, nil
}

func (o *MyOrchestrator) Register(username string, password string, email string, phone string) (bool, error) {
	conn := o.balancer.GetAuth()

	if conn == nil {
		return false, nil
	}

	c := pb.NewAuthenticationClient(conn)
	ctxAlive, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	result, err := c.Register(ctxAlive, &pb.PlayerDetails{Username: username, Password: password, Email: email, Phone: phone})
	if err != nil {
		return false, err
	}

	log.Printf("Register of %s result %t", username, result.Result)
	return result.Result, nil
}
