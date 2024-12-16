package internal

import (
	"log"

	"google.golang.org/grpc"
)

type LoadBalancer interface {
	RegisterService(name string, conn *grpc.ClientConn)
	GetService(name string) *grpc.ClientConn
}

type RandomLoadBalancer struct {
	// TODO: mutex, map with service name and connection object
}

func NewRandomLoadBalancer() *RandomLoadBalancer {
	return &RandomLoadBalancer{}
}

func (lb *RandomLoadBalancer) RegisterService(name string, conn *grpc.ClientConn) {
	log.Printf("SERVICE REGISTERED")
}

func (lb *RandomLoadBalancer) GetService(name string) *grpc.ClientConn {
	return nil
}
