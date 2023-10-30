package main

import (
	"chittychat/proto"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

var clients map[int]proto.MessageService_MessageRouteServer
var idCounter int

func main() {
	clients = make(map[int]proto.MessageService_MessageRouteServer)
	println("Starting Serer")
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
	log.Printf(("New client connected"))
	connID := idCounter
	idCounter++
	clients[connID] = stream
	defer delete(clients, connID)

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Printf("Got message %s, author: %s ", in.Text, in.Author)
		for id, client := range clients {
			if id == connID {
				continue
			}
			if err := client.Send(in); err != nil {
				return err
			}
		}
	}
}
