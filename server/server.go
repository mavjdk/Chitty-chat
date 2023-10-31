package main

import (
	"chittychat/proto"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

var clients map[int32]proto.MessageService_MessageRouteServer
var idCounter int32 = 0
var vectorClock []int32

func main() {
	//make a log file
	f, err := os.OpenFile("../logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	//init vector clock with server position as 0
	vectorClock = []int32{0}

	//print the vector clock
	log.Printf("Server: Starting VectorClock: %v", vectorClock)

	clients = make(map[int32]proto.MessageService_MessageRouteServer)
	grpcServer := grpc.NewServer()
	proto.RegisterMessageServiceServer(grpcServer, &MessageServiceServer{})
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
}

func (s *MessageServiceServer) MessageRoute(stream proto.MessageService_MessageRouteServer) error {
	connID := idCounter
	clients[connID] = stream
	defer removeClient(connID)
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
		log.Printf("Server: Got message %s, author: %s ", in.Text, in.Author)
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

func (s *MessageServiceServer) AddClient(ctx context.Context, in *proto.Empty) (*proto.AddClientResponse, error) {
	vectorClock = append(vectorClock, 1)
	//update because we are receiving a ask from a client
	updateLocalVectorClock()
	var id = updateAndGetId()
	log.Printf("Server: Added client with id: %d", id)
	//update because we are sending the vector clock to the client
	updateLocalVectorClock()
	log.Printf("Server: Sending vectorclock to client")
	return &proto.AddClientResponse{
		Id:          id,
		VectorClock: vectorClock,
	}, nil
}

func removeClient(leavingID int32) error {
	log.Printf("Client %d leaves", leavingID)
	delete(clients, leavingID)
	for id, client := range clients {
		updateLocalVectorClock()
		msg := &proto.Message{
			Text:        fmt.Sprintf("Client %d leaves at vector clock %v", leavingID, vectorClock),
			VectorClock: vectorClock,
			Author:      "Server",
			Id:          id,
		}
		if err := client.Send(msg); err != nil {
			return err
		}
	}
	return nil
}

func updateAndGetId() int32 {
	//send vector clock to the client that requested it
	idCounter++
	return idCounter
}

func updateVectorClockFromClient(clock []int32, id int32) {
	updateLocalVectorClock()
	//because the server is the only one who speaks to clients
	//we only need to update the space for the client in the vector clock
	vectorClock[id] = clock[id]
	//print the vector clock
	log.Printf("Server: Updated vectorclock from client VectorClock: %v", vectorClock)
}
func updateLocalVectorClock() {
	vectorClock[0]++
	log.Printf("Server: Update local vectorclock VectorClock: %v", vectorClock)
}
