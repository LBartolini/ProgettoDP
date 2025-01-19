package services

import (
	"context"
	pb "orchestrator/proto"
	"time"

	"google.golang.org/grpc"
)

type StillAlive interface {
	StillAlive() bool
	Close()
}

func StillAliveHandle(conn *grpc.ClientConn) bool {
	// Function used by service to check if underlying connection is still alive

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := pb.NewStillAliveClient(conn).StillAlive(ctx, nil)

	return err == nil
}
