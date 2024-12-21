package main

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"fmt"
	"log"
	"net"
	"os"
	"time"

	"garage/internal"

	pb "garage/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("SERVICE_PORT")))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("server listening at %v", lis.Addr())

	db, err := sql.Open("mysql", "root:admin@tcp(garage_db:3306)/Garage")
	if err != nil {
		log.Fatalf("failed to connect to db: %s", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("error pinging database: %v", err)
	}

	s := grpc.NewServer()
	server := internal.NewServer(internal.NewSQL_DB(db))
	pb.RegisterGarageServer(s, server)
	pb.RegisterStillAliveServer(s, server)

	go registerToOrchestrator()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

func registerToOrchestrator() {
	log.Printf("Trying to connect to Orchestrator")
	conn, err := grpc.NewClient(fmt.Sprintf("orchestrator:%s", os.Getenv("SERVICE_PORT")), grpc.WithTransportCredentials(insecure.NewCredentials()))
	for err != nil {
		log.Print(err.Error())
		time.Sleep(500 * time.Millisecond)
		conn, err = grpc.NewClient(fmt.Sprintf("orchestrator:%s", os.Getenv("SERVICE_PORT")), grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	defer conn.Close()

	c := pb.NewOrchestratorClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = c.RegisterGarage(ctx, nil)
	for err != nil {
		log.Print(err.Error())
		time.Sleep(500 * time.Millisecond)
		_, err = c.RegisterGarage(ctx, nil)
	}
	log.Printf("Registered to Orchestrator")
}
