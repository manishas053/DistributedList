package main

import (
	"log"
	"time"
	"fmt"
	"bytes"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"github.com/Nik-U/pbc"
	pb "github.com/DistributedList/list"
)

const (
	address     = "10.1.4.3:50051"
	defaultName = "list"
)

func adding(client pb.ListClient, in *pb.InputMsg) {
	start := time.Now()
	var temp []byte
	resp, err := client.ProcessInput(context.Background(), in)
	if err != nil {
		log.Fatalf("node couldn't insert %v", err)
	}
	if !bytes.Equal(resp.Resp, temp) {
		log.Printf("BLS Multisignature aggregate at the root is :\n %v", resp.Resp)
	}
	fmt.Println("Time taken to compute BLS multisignature : ", time.Since(start))
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
