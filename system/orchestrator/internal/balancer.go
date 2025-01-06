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

type LoadBalancer interface {
	RegisterAuth(conn *grpc.ClientConn) error
	GetAuth() *grpc.ClientConn

	RegisterLeaderboard(conn *grpc.ClientConn) error
	GetLeaderboard() *grpc.ClientConn

	RegisterGarage(conn *grpc.ClientConn) error
	GetGarage() *grpc.ClientConn

	RegisterRacing(conn *grpc.ClientConn) error
	GetRacing() *grpc.ClientConn

	testConnection(*grpc.ClientConn) error
}

type RandomLoadBalancer struct {
	mu       sync.Mutex
	services map[string][]*grpc.ClientConn
}

func NewRandomLoadBalancer() *RandomLoadBalancer {
	return &RandomLoadBalancer{services: make(map[string][]*grpc.ClientConn)}
}

func (lb *RandomLoadBalancer) registerService(name string, conn *grpc.ClientConn) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if _, exists := lb.services[name]; !exists {
		lb.services[name] = []*grpc.ClientConn{}
	}
	lb.services[name] = append(lb.services[name], conn)
	return nil
}

func removeAtIndex(slice []*grpc.ClientConn, index int) []*grpc.ClientConn {
	if index < 0 || index >= len(slice) {
		return slice
	}
	return append(slice[:index], slice[index+1:]...)
}

func (lb *RandomLoadBalancer) getService(name string) *grpc.ClientConn {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	var conn *grpc.ClientConn

	for i := len(lb.services[name]); i > 0; i-- {
		index := rand.IntN(len(lb.services[name]))
		temp := lb.services[name][index]

		if err := lb.testConnection(temp); err == nil {
			conn = temp
			break
		} else {
			log.Printf("Balancer removing service %s not alive", name)
			lb.services[name] = removeAtIndex(lb.services[name], index)
			temp.Close()
		}
	}

	return conn
}

func (lb *RandomLoadBalancer) RegisterAuth(conn *grpc.ClientConn) error {
	return lb.registerService("auth", conn)
}

func (lb *RandomLoadBalancer) GetAuth() *grpc.ClientConn {
	return lb.getService("auth")
}

func (lb *RandomLoadBalancer) RegisterLeaderboard(conn *grpc.ClientConn) error {
	return lb.registerService("leaderboard", conn)
}

func (lb *RandomLoadBalancer) GetLeaderboard() *grpc.ClientConn {
	return lb.getService("leaderboard")
}

func (lb *RandomLoadBalancer) RegisterRacing(conn *grpc.ClientConn) error {
	return lb.registerService("racing", conn)
}

func (lb *RandomLoadBalancer) GetRacing() *grpc.ClientConn {
	return lb.getService("racing")
}

func (lb *RandomLoadBalancer) RegisterGarage(conn *grpc.ClientConn) error {
	return lb.registerService("garage", conn)
}

func (lb *RandomLoadBalancer) GetGarage() *grpc.ClientConn {
	return lb.getService("garage")
}

func (lb *RandomLoadBalancer) testConnection(conn *grpc.ClientConn) error {
	c := pb.NewStillAliveClient(conn)
	ctxAlive, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := c.StillAlive(ctxAlive, nil)
	return err
}
