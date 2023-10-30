package main

import (
	"context"
	"io"
	"log"
	"sync"

	"chittychat/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var waitgroup sync.WaitGroup

func main() {

	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	log.Printf("Connection State: %s", conn.GetState().String())
	defer conn.Close()

	ServiceConn := proto.NewMessageServiceClient(conn)

	stream, err := ServiceConn.MessageRoute(context.Background())
	defer stream.CloseSend()
	if err != nil {
		log.Fatalf("Error when calling MessageRoute: %s", err)
	}
	waitgroup.Add(1)
	go listener(stream)
	msg := &proto.Message{
		Text:   "Hello from the client!",
		Author: "Client",
	}
	if err := stream.Send(msg); err != nil {
		log.Fatalf("Failed to send a msg: %v", err)
	}

	waitgroup.Wait()

}

func listener(stream proto.MessageService_MessageRouteClient) {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			// read done.
			log.Printf("EOF")
			waitgroup.Done()
			return
		}
		if err != nil {
			log.Fatalf("Failed to receive a msg : %v", err)
		}
		log.Printf("Got message %s, author: %s ", in.Text, in.Author)
	}
}
