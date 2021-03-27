package client

import (
	"context"
	"fiy/pkg/grpc/proto/host"
	"log"
	"time"

	"google.golang.org/grpc"
)

/*
  @Author : lanyulei
*/

const (
	address = "localhost:50051"
)

func RunClient() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return
	}
	defer conn.Close()
	c := host.NewHostInfoClient(conn)
	// Contact the server and print out its response.
	data, err := IntegrateData()
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.GetHostInfo(ctx, &host.HostInfoRequest{Data: data})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
}
