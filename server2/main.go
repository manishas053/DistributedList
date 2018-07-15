package main

import (
	"log"
	"net"
	"sync"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"github.com/Nik-U/pbc"
	"crypto/sha256"
	pb "github.com/DistributedList/list"
)

const (
	port       = ":50052"
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
	log.Printf("Signature %v", signature)

	resp, err := L.ProcessInput(context.Background(), &pb.InputMsg{SharedParams : in.SharedParams, SharedG : in.SharedG})
	if err != nil {
		log.Printf("Error connecting %v", err)
	}
	response := resp.Resp
	return &pb.InputResponse{Resp: response}, nil
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
			resp, err := multiserver.InsertInput(context.Background(), &pb.InputMsg{SharedParams: in.SharedParams, SharedG: in.SharedG})
			if err != nil {
				log.Printf("Error connecting %v", err)
			}
			out <- resp
			signature := pairing.NewG2().SetBytes(resp.Resp)
			aggregate = pairing.NewG2().Add(aggregate, signature)
		}()
	}
	i := 1
	for n := range out {
		fmt.Println(n)
		i = i + 1
		if i >= len(L.replicas) + 1 {
			close(out)
			//log.Printf("Aggregate is : %v", aggregate)
		}
	}
	aggregate_sum := aggregate.Bytes()
	return &pb.InputResponse{Resp: aggregate_sum}, nil
}


func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	l := List{}
	l.id = 1
	l.address = "localhost:50052"
	l.replicas = append(l.replicas, "localhost:50055", "localhost:50056")
	s := grpc.NewServer()
	pb.RegisterListServer(s, &l)
	s.Serve(lis)
}
