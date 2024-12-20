package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"orchestrator/internal"

	pb "orchestrator/proto"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

	r := gin.Default()
	r.Use(sessions.Sessions("session", store))
	r.LoadHTMLGlob("./templates/*")

	orchestrator := internal.NewOrchestrator(internal.NewRandomLoadBalancer())
	go startOrchestratorService(orchestrator)
	routes := internal.NewMyRoutes(orchestrator)

	r.GET("/", routes.IndexRoute)

	r.POST("/login", routes.LoginRoute)
	r.POST("/register", routes.RegisterRoute)
	r.GET("/leaderboard", routes.LeaderboardRoute)

	private := r.Group("/private")
	private.Use(internal.Authorized)
	{
		private.GET("/", routes.HomeRoute)
		private.POST("/logout", routes.LogoutRoute)
		private.GET("/garage", routes.GarageRoute)
		private.GET("/history", routes.RaceHistoryRoute)
	}

	r.Run("0.0.0.0:8080")
}

func startOrchestratorService(orchestrator *internal.Orchestrator) {
	log.Printf("Starting Orchestrator Service")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("SERVICE_PORT")))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterOrchestratorServer(s, orchestrator)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
