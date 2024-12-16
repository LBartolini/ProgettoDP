/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	pb "auth/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedAuthenticationServer
	pb.UnimplementedStillAliveServer
}

func (s *server) Login(ctx context.Context, in *pb.PlayerCredentials) (*pb.LoginResult, error) {
	return &pb.LoginResult{Result: true}, nil
}

func (s *server) Register(ctx context.Context, in *pb.PlayerDetails) (*pb.RegisterResult, error) {
	return &pb.RegisterResult{Result: true, Id: 1}, nil
}

func (s *server) StillAlive(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	log.Printf("ALIVE")
	return nil, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("SERVICE_PORT")))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAuthenticationServer(s, &server{})
	pb.RegisterStillAliveServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
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
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = c.RegisterService(ctx, &pb.RegisterServiceMessage{Name: "auth"})
	for err != nil {
		log.Print(err.Error())
		time.Sleep(500 * time.Millisecond)
		_, err = c.RegisterService(ctx, &pb.RegisterServiceMessage{Name: "auth"})
	}
	log.Printf("Registered to Orchestrator")
}
