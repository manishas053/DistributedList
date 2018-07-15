package main

import (
	"log"
	"bytes"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"github.com/Nik-U/pbc"
	pb "github.com/DistributedList/list"
)

const (
	address     = "localhost:50051"
	defaultName = "list"
)

func adding(client pb.ListClient, in *pb.InputMsg) {
	var temp []byte
	resp, err := client.ProcessInput(context.Background(), in)
	if err != nil {
		log.Fatalf("node couldn't insert %v", err)
	}
	if !bytes.Equal(resp.Resp, temp) {
		log.Printf("Aggregate at the root is : %v", resp.Resp)
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


	params := pbc.GenerateA(160, 512)
  pairing := params.NewPairing()
 	g := pairing.NewG2().Rand()

	in := &pb.InputMsg{SharedParams : params.String(), SharedG : g.Bytes()}

	adding(client, in)
}
