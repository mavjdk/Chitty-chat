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
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	log.Printf("Connection State: %s", conn.GetState().String())
	defer conn.Close()

	initialVectorClockAsk(proto.NewVectorClockServiceClient(conn))

	ServiceConn := proto.NewMessageServiceClient(conn)

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
		msg := &proto.Message{
			Text:        input,
			VectorClock: vectorClock,
			Id:          id,
			Author:      fmt.Sprintf("Client %d", id),
		}
		updateLocalVectorClock()
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
		updateVectorClock(in.VectorClock, in.Id)
		log.Printf("Got message %s, author: %s", in.Text, in.Author)
		//print the vector clock
		log.Printf("VectorClock: %v", vectorClock)
	}
}
func initialVectorClockAsk(client proto.VectorClockServiceClient) {
	respone, err := client.VectorClockService(context.Background(), &proto.Empty{})
	if err != nil {
		log.Fatalf("Failed to get VectorClock: %v", err)
	}
	vectorClock = respone.VectorClock
	id = respone.Id

}
func updateLocalVectorClock() {
	vectorClock[id]++
}
func updateVectorClock(incommingClock []int32, id int32) {
	if vectorClock[id] < incommingClock[id] {
		vectorClock[id] = incommingClock[id]
	}
}
