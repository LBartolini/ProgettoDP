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

func (o *MyOrchestrator) RegisterAuth(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	peerInfo, ok := peer.FromContext(ctx)
	if !ok {
		return nil, errors.New("unable to get peer information")
	}

	address := strings.Split(peerInfo.Addr.String(), ":")[0]
	log.Printf("%s", address)
	client, err := grpc.NewClient(fmt.Sprintf("%s:%s", address, os.Getenv("SERVICE_PORT")), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		client.Close()
		return nil, err
	}

	return nil, o.balancer.RegisterAuth(client)
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
		log.Printf("error during login")
		return false, err
	}

	log.Printf("Login: %t", result.Result)
	return result.Result, nil
}
