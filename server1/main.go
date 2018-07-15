package main

import (
	"fmt"
	"log"
	"net"
	"sync"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"github.com/Nik-U/pbc"
	"crypto/sha256"
	pb "github.com/DistributedList/list"
)

const (
	port       = ":50051"
	reply_port = ":50052"
)

// Represent a node
type Node struct {
	data int32
	next *Node
	prev *Node
}

// Represent a linked list
type List struct {
	id       int32
	address  string
	replicas []string

	head *Node
	tail *Node

	reply chan *pb.InputResponse
}

func (L *List) InsertInput(ctx context.Context, in *pb.InputMsg) (*pb.InputResponse, error) {
	params, _ := pbc.NewPairingFromString(in.SharedParams)
	privateKey := params.NewZr().Rand()
	message := "some text to sign"
	h := params.NewG1().SetFromStringHash(message, sha256.New())
	sign := params.NewG2().PowZn(h, privateKey)
	signature := sign.Bytes()
	log.Printf("signature %v", signature)


	//response := int32(value)
	return &pb.InputResponse{Resp: signature}, nil
}

//Function to do multicast to all the replicas
func (L *List) ProcessInput(ctx context.Context, in *pb.InputMsg) (*pb.InputResponse, error) {
	var wg sync.WaitGroup
	pairing, _ := pbc.NewPairingFromString(in.SharedParams)
	aggregate := pairing.NewG2()
	out := make(chan *pb.InputResponse)
	wg.Add(len(L.replicas))
	for _, replica := range L.replicas {
		cliConn, err := grpc.Dial(replica, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Didn't connect %v", err)
		}
		multiserver := pb.NewListClient(cliConn)
		go func() {
			resp, err := multiserver.InsertInput(context.Background(), &pb.InputMsg{SharedParams : in.SharedParams, SharedG : in.SharedG})
			if err != nil {
				log.Printf("Error connecting %v", err)
			}
			out <- resp
			signature := pairing.NewG2().SetBytes(resp.Resp)
			aggregate = pairing.NewG2().Add(aggregate, signature)
		}()
	}
	i := 0
	for n := range out {
		fmt.Println(n)
		i = i + 1
		if i >= len(L.replicas) {
			close(out)
		}
	}
	aggregate_sum := aggregate.Bytes()
	log.Printf("Aggregate %v", aggregate_sum)
	return &pb.InputResponse{Resp: aggregate_sum}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	l := List{}
	l.id = 1
	l.address = "localhost:50051"
	l.replicas = append(l.replicas, "localhost:50052", "localhost:50050")
	s := grpc.NewServer()
	pb.RegisterListServer(s, &l)
	s.Serve(lis)
}
