package main

import (
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"github.com/Nik-U/pbc"
	"crypto/sha256"
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
	params, _ := pbc.NewPairingFromString(in.SharedParams)
	privateKey := params.NewZr().Rand()
	message := "some text to sign"
  h := params.NewG1().SetFromStringHash(message, sha256.New())
	sign := params.NewG2().PowZn(h, privateKey)
	signature := sign.Bytes()
	log.Printf("Signature %v", signature)

	return &pb.InputResponse{Resp: signature}, nil
}

//Function to do multicast to all the replicas
func (L *List) ProcessInput(ctx context.Context, in *pb.InputMsg) (*pb.InputResponse, error) {
	const (
		address = "localhost:50053"
		)
	//for _, replica := range L.replicas {
		cliConn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Didn't connect %v", err)
		}
		multiserver := pb.NewListClient(cliConn)
		resp, err := multiserver.InsertInput(context.Background(), &pb.InputMsg{SharedParams : in.SharedParams, SharedG : in.SharedG})
		if err != nil {
			log.Printf("Error connecting %v", err)
		}
		log.Print(resp)
		//value = value + resp
		response := resp.Resp
	//}
	return &pb.InputResponse{Resp: response}, nil
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
