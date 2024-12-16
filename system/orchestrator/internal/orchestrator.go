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

type Orchestrator interface {
}

type MyOrchestrator struct {
	pb.UnimplementedOrchestratorServer
	balancer LoadBalancer
}

func NewMyOrchestrator(balancer LoadBalancer) *MyOrchestrator {
	return &MyOrchestrator{balancer: balancer}
}

func (o *MyOrchestrator) RegisterService(ctx context.Context, in *pb.RegisterServiceMessage) (*emptypb.Empty, error) {
	peerInfo, ok := peer.FromContext(ctx)
	if !ok {
		log.Printf("Unable to get peer information")
		return nil, errors.New("unable to get peer information")
	}

	address := strings.Split(peerInfo.Addr.String(), ":")[0]
	log.Printf("%s", address)
	client, err := grpc.NewClient(fmt.Sprintf("%s:%s", address, os.Getenv("SERVICE_PORT")), grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer client.Close()
	if err != nil {
		log.Printf("Unable to connect to service")
		return nil, err
	}

	c := pb.NewStillAliveClient(client)
	ctxAlive, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = c.StillAlive(ctxAlive, nil)
	if err != nil {
		log.Printf("Unable to connect to service")
		log.Print(err.Error())
		return nil, err
	}

	o.balancer.RegisterService(in.Name, client)
	return nil, nil
}
