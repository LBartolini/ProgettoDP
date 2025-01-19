package internal

import (
	"context"
	"log"
	pb "orchestrator/proto"
	"sync"
	"time"

	"math/rand/v2"

	"google.golang.org/grpc"
)

// TODO: balancer works with services, orchestrator no longer has to do the work of creating the service from connection
type LoadBalancer interface {
	RegisterAuth(conn *grpc.ClientConn)
	GetAuth() *grpc.ClientConn

	RegisterLeaderboard(conn *grpc.ClientConn)
	GetLeaderboard() *grpc.ClientConn

	RegisterGarage(conn *grpc.ClientConn)
	GetGarage() *grpc.ClientConn

	RegisterRacing(conn *grpc.ClientConn)
	GetRacing() *grpc.ClientConn

	testConnection(*grpc.ClientConn) error
}

// Implementation of LoadBalancer with random selection
type RandomLoadBalancer struct {
	mu       sync.Mutex
	services map[string][]*grpc.ClientConn
}

func NewRandomLoadBalancer() *RandomLoadBalancer {
	return &RandomLoadBalancer{services: make(map[string][]*grpc.ClientConn)}
}

func (lb *RandomLoadBalancer) registerService(name string, conn *grpc.ClientConn) {
	// Register service, if it is the first time also create the underlying slice

	lb.mu.Lock()
	defer lb.mu.Unlock()

	if _, exists := lb.services[name]; !exists {
		lb.services[name] = []*grpc.ClientConn{}
	}
	lb.services[name] = append(lb.services[name], conn)
}

func removeAtIndex(slice []*grpc.ClientConn, index int) []*grpc.ClientConn {
	// Utility function to remove connection at specified index

	if index < 0 || index >= len(slice) {
		return slice
	}
	return append(slice[:index], slice[index+1:]...)
}

func (lb *RandomLoadBalancer) getService(name string) *grpc.ClientConn {
	// Retrieve random replica of service

	lb.mu.Lock()
	defer lb.mu.Unlock()
	var conn *grpc.ClientConn

	for i := len(lb.services[name]); i > 0; i-- {
		index := rand.IntN(len(lb.services[name]))
		temp := lb.services[name][index]

		// test connection using StillAlive service
		if err := lb.testConnection(temp); err == nil {
			conn = temp
			break
		} else {
			// Selected replica not alive, removing and retrying
			log.Printf("Balancer removing service, %s not alive", name)
			lb.services[name] = removeAtIndex(lb.services[name], index)
			temp.Close()
		}
	}

	return conn
}

func (lb *RandomLoadBalancer) RegisterAuth(conn *grpc.ClientConn) {
	lb.registerService("auth", conn)
}

func (lb *RandomLoadBalancer) GetAuth() *grpc.ClientConn {
	return lb.getService("auth")
}

func (lb *RandomLoadBalancer) RegisterLeaderboard(conn *grpc.ClientConn) {
	lb.registerService("leaderboard", conn)
}

func (lb *RandomLoadBalancer) GetLeaderboard() *grpc.ClientConn {
	return lb.getService("leaderboard")
}

func (lb *RandomLoadBalancer) RegisterRacing(conn *grpc.ClientConn) {
	lb.registerService("racing", conn)
}

func (lb *RandomLoadBalancer) GetRacing() *grpc.ClientConn {
	return lb.getService("racing")
}

func (lb *RandomLoadBalancer) RegisterGarage(conn *grpc.ClientConn) {
	lb.registerService("garage", conn)
}

func (lb *RandomLoadBalancer) GetGarage() *grpc.ClientConn {
	return lb.getService("garage")
}

func (lb *RandomLoadBalancer) testConnection(conn *grpc.ClientConn) error {
	c := pb.NewStillAliveClient(conn)
	ctxAlive, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// if replica is not alive an error will be raised
	_, err := c.StillAlive(ctxAlive, nil)
	return err
}
