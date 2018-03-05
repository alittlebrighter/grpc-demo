//go:generate protoc -I defns/ --go_out=plugins=grpc:defns defns/service.proto defns/models.proto
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "github.com/alittlebrighter/grpc-demo/defns"
)

func main() {
	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial(":4444", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewSampleClient(conn)

	greeting, err := client.Greet(context.Background(), &pb.Greetee{FirstName: "Vinny", LastName: "Veriwolf"})
	if err != nil {
		log.Println("was not greeted: " + err.Error())
	}
	fmt.Println("Greeting received: " + greeting.GetGreeting())

	vals := []*pb.Val{
		&pb.Val{Label: "one", Value: 235},
		&pb.Val{Label: "two", Value: -764},
		&pb.Val{Label: "three", Value: 1024},
	}
	intervals := []time.Duration{
		500 * time.Millisecond,
		1000 * time.Millisecond,
		1500 * time.Millisecond,
	}

	// streams are different in that each request/response cycle is a new struct
	lifo, err := client.LifoEcho(context.Background())
	if err != nil {
		log.Fatalln("was not greeted: " + err.Error())
	}

	for i, val := range vals {
		time.Sleep(intervals[i])
		fmt.Printf("Sending: %s %d\n", val.GetLabel(), val.GetValue())
		lifo.Send(val)
	}
	lifo.CloseSend()

	for {
		val, err := lifo.Recv()
		if err == io.EOF {
			fmt.Println("Response stream closed.")
			break
		}

		fmt.Printf("Received: %s %d\n", val.GetLabel(), val.GetValue())
	}
}
