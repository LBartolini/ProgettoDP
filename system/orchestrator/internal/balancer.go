package internal

import (
	"orchestrator/internal/services"
	"sync"

	"math/rand/v2"
)

type LoadBalancer interface {
	RegisterAuth(services.Auth)
	GetAuth() services.Auth

	RegisterLeaderboard(services.Leaderboard)
	GetLeaderboard() services.Leaderboard

	RegisterGarage(services.Garage)
	GetGarage() services.Garage

	RegisterRacing(services.Racing)
	GetRacing() services.Racing
}

// Implementation of LoadBalancer with random selection
type RandomLoadBalancer struct {
	mu_auth sync.Mutex
	auth    []services.Auth

	mu_leaderboard sync.Mutex
	leaderboard    []services.Leaderboard

	mu_garage sync.Mutex
	garage    []services.Garage

	mu_racing sync.Mutex
	racing    []services.Racing
}

func NewRandomLoadBalancer() *RandomLoadBalancer {
	return &RandomLoadBalancer{}
}

func (lb *RandomLoadBalancer) RegisterAuth(s services.Auth) {
	lb.mu_auth.Lock()
	defer lb.mu_auth.Unlock()

	lb.auth = append(lb.auth, s)
}

func (lb *RandomLoadBalancer) GetAuth() (s services.Auth) {
	lb.mu_auth.Lock()
	defer lb.mu_auth.Unlock()

	for i := len(lb.auth); i > 0; i-- {
		index := rand.IntN(len(lb.auth))
		temp := lb.auth[index]

		// test connection using StillAlive service
		if temp.StillAlive() {
			s = temp
			break
		} else {
			// Selected replica not alive, removing and retrying
			lb.auth = append(lb.auth[:index], lb.auth[index+1:]...)
			temp.Close()
		}
	}

	return s
}

func (lb *RandomLoadBalancer) RegisterLeaderboard(s services.Leaderboard) {
	lb.mu_leaderboard.Lock()
	defer lb.mu_leaderboard.Unlock()

	lb.leaderboard = append(lb.leaderboard, s)
}

func (lb *RandomLoadBalancer) GetLeaderboard() (s services.Leaderboard) {
	lb.mu_leaderboard.Lock()
	defer lb.mu_leaderboard.Unlock()

	for i := len(lb.leaderboard); i > 0; i-- {
		index := rand.IntN(len(lb.leaderboard))
		temp := lb.leaderboard[index]

		// test connection using StillAlive service
		if temp.StillAlive() {
			s = temp
			break
		} else {
			// Selected replica not alive, removing and retrying
			lb.leaderboard = append(lb.leaderboard[:index], lb.leaderboard[index+1:]...)
			temp.Close()
		}
	}

	return s
}

func (lb *RandomLoadBalancer) RegisterRacing(s services.Racing) {
	lb.mu_racing.Lock()
	defer lb.mu_racing.Unlock()

	lb.racing = append(lb.racing, s)
}

func (lb *RandomLoadBalancer) GetRacing() (s services.Racing) {
	lb.mu_racing.Lock()
	defer lb.mu_racing.Unlock()

	for i := len(lb.racing); i > 0; i-- {
		index := rand.IntN(len(lb.racing))
		temp := lb.racing[index]

		// test connection using StillAlive service
		if temp.StillAlive() {
			s = temp
			break
		} else {
			// Selected replica not alive, removing and retrying
			lb.racing = append(lb.racing[:index], lb.racing[index+1:]...)
			temp.Close()
		}
	}

	return s
}

func (lb *RandomLoadBalancer) RegisterGarage(s services.Garage) {
	lb.mu_garage.Lock()
	defer lb.mu_garage.Unlock()

	lb.garage = append(lb.garage, s)
}

func (lb *RandomLoadBalancer) GetGarage() (s services.Garage) {
	lb.mu_garage.Lock()
	defer lb.mu_garage.Unlock()

	for i := len(lb.garage); i > 0; i-- {
		index := rand.IntN(len(lb.garage))
		temp := lb.garage[index]

		// test connection using StillAlive service
		if temp.StillAlive() {
			s = temp
			break
		} else {
			// Selected replica not alive, removing and retrying
			lb.garage = append(lb.garage[:index], lb.garage[index+1:]...)
			temp.Close()
		}
	}

	return s
}
