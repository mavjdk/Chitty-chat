package main

import (
	"bufio"
	"chittychat/proto"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var vectorClock []int32
var id int32 = 1

func main() {
	//make a log file
	f, err := os.OpenFile("../logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	log.Printf("Connecting to server, State: %s", conn.GetState().String())
	defer conn.Close()

	ServiceConn := proto.NewMessageServiceClient(conn)

	initialAddClientCall(ServiceConn)

	stream, err := ServiceConn.MessageRoute(context.Background())
	defer stream.CloseSend()
	if err != nil {
		log.Fatalf("Error when calling MessageRoute: %s", err)
	}
	go listener(stream)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		input := scanner.Text()
		if err != nil {
			log.Fatalf("Failed to scan input: %v", err)
		}
		log.Printf("Client %d: Sending message", id)
		updateLocalVectorClock()
		msg := &proto.Message{
			Text:        input,
			VectorClock: vectorClock,
			Id:          id,
			Author:      fmt.Sprintf("Client %d", id),
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
		updateVectorClock(in.VectorClock)
		fmt.Printf("Got message %s, author: %s\n", in.Text, in.Author)
		log.Printf("Got message %s, author: %s", in.Text, in.Author)
	}
}
func initialAddClientCall(client proto.MessageServiceClient) {
	response, err := client.AddClient(context.Background(), &proto.Empty{})
	if err != nil {
		log.Fatalf("Failed to Add client to server: %v", err)
	}
	vectorClock = response.VectorClock
	id = response.Id
	//update because we are receiving
	updateLocalVectorClock()

}
func updateLocalVectorClock() {
	vectorClock[id]++
	log.Printf("Clint %d: Update local vectorclock VectorClock: %v", id, vectorClock)
}
func updateVectorClock(incommingClock []int32) {
	vectorClock[id]++
	if len(vectorClock) < len(incommingClock) {
		vectorClock = append(vectorClock, incommingClock[len(vectorClock):]...)
	}
	for i := 0; i < len(vectorClock); i++ {
		if vectorClock[i] < incommingClock[i] {
			vectorClock[i] = incommingClock[i]
		}
	}
	log.Printf("Client %d: Update vectorclock based on message from server VectorClock: %v", id, vectorClock)
}
