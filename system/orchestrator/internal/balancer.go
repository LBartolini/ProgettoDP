package internal

import (
	"context"
	"log"
	pb "orchestrator/proto"
	"time"

	"google.golang.org/grpc"
)

type LoadBalancer interface {
	RegisterService(name string, conn *grpc.ClientConn) error
	GetService(name string) *grpc.ClientConn
	testConnection(*grpc.ClientConn) error
}

type RandomLoadBalancer struct {
	// TODO: mutex, map with service name and connection object
}

func NewRandomLoadBalancer() *RandomLoadBalancer {
	return &RandomLoadBalancer{}
}

func (lb *RandomLoadBalancer) RegisterService(name string, conn *grpc.ClientConn) error {
	if err := lb.testConnection(conn); err != nil {
		return err
	}

	log.Printf("SERVICE REGISTERED")
	return nil
}

func (lb *RandomLoadBalancer) GetService(name string) *grpc.ClientConn {
	return nil
}

func (lb *RandomLoadBalancer) testConnection(conn *grpc.ClientConn) error {
	c := pb.NewStillAliveClient(conn)
	ctxAlive, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := c.StillAlive(ctxAlive, nil)

	return err
}
