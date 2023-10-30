package main

import (
	"chittychat/proto"
	"context"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

var clients map[int32]proto.MessageService_MessageRouteServer
var idCounter int32 = 0
var vectorClock []int32

func main() {
	//init vector clock with server position as 0
	vectorClock = make([]int32, 10)
	//print the vector clock
	log.Printf("VectorClock: %v", vectorClock)

	clients = make(map[int32]proto.MessageService_MessageRouteServer)
	println("Starting Serer")
	grpcServer := grpc.NewServer()
	proto.RegisterMessageServiceServer(grpcServer, &MessageServiceServer{})
	proto.RegisterVectorClockServiceServer(grpcServer, &MessageServiceServer{})
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}

type MessageServiceServer struct {
	proto.UnimplementedMessageServiceServer
	proto.UnimplementedVectorClockServiceServer
}

func (s *MessageServiceServer) MessageRoute(stream proto.MessageService_MessageRouteServer) error {
	log.Printf(("New client connected"))
	connID := idCounter
	clients[connID] = stream
	defer delete(clients, connID)

	for {
		//receive message from client
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		updateVectorClockFromClient(in.VectorClock, in.Id)
		log.Printf("Got message %s, author: %s ", in.Text, in.Author)
		for id, client := range clients {
			//skips sending to the client that sent the message
			if id == connID {
				continue
			}
			updateLocalVectorClock()
			in.VectorClock = vectorClock
			if err := client.Send(in); err != nil {
				return err
			}
		}
	}
}

func (s *MessageServiceServer) VectorClockService(ctx context.Context, in *proto.Empty) (*proto.VectorClock, error) {
	updateLocalVectorClock()
	var id = updateAndGetId()
	vectorClock[id] = 0

	return &proto.VectorClock{
		Id:          id,
		VectorClock: vectorClock,
	}, nil
}

func updateAndGetId() int32 {
	//send vector clock to the client that requested it
	idCounter++
	return idCounter
}

func updateVectorClockFromClient(clock []int32, id int32) {
	//print the vector clock
	log.Printf("Got VectorClock From client: %v", clock)
	//print the id of the client that sent the vector clock

	updateLocalVectorClock()
	vectorClock[0]++
	//because the server is the only one who speaks to clients
	//we only need to update the space for the client in the vector clock
	vectorClock[id] = clock[id]
}
func updateLocalVectorClock() {
	vectorClock[0]++
}
