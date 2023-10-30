package main

import (
	"chittychat/proto"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {

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
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Printf("Got message %s, author: %s ", in.Text, in.Author)
		stream.Send(&proto.Message{
			Text:   "Hello from the server!",
			Author: "Server",
		})

	}
}
