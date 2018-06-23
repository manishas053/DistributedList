package main

import (
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/DistributedList/list"
)

const (
	port       = ":50053"
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
	log.Printf("Input Message %v", in.Data)
	value := 2
	response := int32(value)
	return &pb.InputResponse{Resp: response}, nil
}

//Function to do multicast to all the replicas
func (L *List) ProcessInput(ctx context.Context, in *pb.InputMsg) (*pb.InputResponse, error) {
	//value := 0
	for _, replica := range L.replicas {
		cliConn, err := grpc.Dial(replica, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Didn't connect %v", err)
		}
		multiserver := pb.NewListClient(cliConn)
		resp, err := multiserver.InsertInput(context.Background(), &pb.InputMsg{Data: 2})
		if err != nil {
			log.Printf("Error connecting %v", err)
		}
		log.Print(resp)
		//value = value + resp
	}
	return &pb.InputResponse{Resp: 1}, nil
}


func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	l := List{}
	l.id = 1
	l.address = "localhost:50053"
	l.replicas = append(l.replicas, "localhost:50050", "localhost:50053")
	s := grpc.NewServer()
	pb.RegisterListServer(s, &l)
	s.Serve(lis)
}
