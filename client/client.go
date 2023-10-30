package main

import (
	"bufio"
	"chittychat/proto"
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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
	go listener(stream)
	scanner := bufio.NewScanner(os.Stdin)
	clientID := rand.Int() % 1000
	for {
		scanner.Scan()
		input := scanner.Text()
		if err != nil {
			log.Fatalf("Failed to scan input: %v", err)
		}
		msg := &proto.Message{
			Text:   input,
			Author: fmt.Sprintf("Client %d", clientID),
		}
		if err := stream.Send(msg); err != nil {
			log.Fatalf("Failed to send a msg: %v", err)
		}
	}
}

func listener(stream proto.MessageService_MessageRouteClient) {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			// read done.
			log.Printf("EOF")
			return
		}
		if err != nil {
			log.Fatalf("Failed to receive a msg : %v", err)
		}
		log.Printf("Got message %s, author: %s ", in.Text, in.Author)
	}
}
