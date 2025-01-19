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
	// Cookie store using env variable session_key
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

	r := gin.Default()
	r.Use(sessions.Sessions("session", store))
	r.LoadHTMLGlob("./templates/*") // load templates

	// Create and start orchestrator gRPC service
	orchestrator := internal.NewOrchestrator(internal.NewRandomLoadBalancer())
	go startOrchestratorService(orchestrator)

	// Handle routes
	routes := internal.NewMyRoutes(orchestrator)

	r.GET("/", routes.IndexRoute)

	r.POST("/login", routes.LoginRoute)
	r.POST("/register", routes.RegisterRoute)
	r.GET("/leaderboard", routes.LeaderboardRoute)

	// Group of routes that need Authorization
	private := r.Group("/private")
	private.Use(internal.Authorized)
	{
		private.GET("/", routes.HomeRoute)
		private.POST("/logout", routes.LogoutRoute)

		private.GET("/garage", routes.GarageRoute)
		private.POST("/garage/buy", routes.GarageBuyRoute)
		private.POST("/garage/upgrade", routes.GarageUpgradeRoute)

		private.POST("/race/start", routes.RaceStartRoute)

		private.GET("/history", routes.RaceHistoryRoute)
	}

	// Run webserver on env web_port
	r.Run(fmt.Sprintf("0.0.0.0:%s", os.Getenv("WEB_PORT")))
}

func startOrchestratorService(orchestrator *internal.Orchestrator) {
	log.Printf("Starting Orchestrator Service")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("SERVICE_PORT")))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Creation of gRPC Server
	s := grpc.NewServer()
	pb.RegisterOrchestratorServer(s, orchestrator)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
