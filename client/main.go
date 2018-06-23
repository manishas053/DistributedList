package main

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/DistributedList/list"
)

const (
	address     = "localhost:50051"
	defaultName = "list"
)

func adding(client pb.ListClient, in *pb.InputMsg) {
	resp, err := client.ProcessInput(context.Background(), in)
	if err != nil {
		log.Fatalf("node couldn't insert %v", err)
	}
	if resp.Resp != 0 {
		log.Printf("New Node inserted %v", in.Data)
	}
}

func main() {
	// Setup connection to gRPC server
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Didn't connect %v", err)
	}
	defer conn.Close()
	client := pb.NewListClient(conn)

	in := &pb.InputMsg{Data: 1}
	adding(client, in)
}
