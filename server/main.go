package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	pb "github.com/alittlebrighter/grpc-demo/defns"
)

func main() {
	sampleSvc := new(SampleServer)
	listen, err := net.Listen("tcp", ":4444")
	if err != nil {
		log.Fatalf("failed to listen for tcp connections: %v", err)
	}

	rpcServer := grpc.NewServer()
	pb.RegisterSampleServer(rpcServer, sampleSvc)
	if err := rpcServer.Serve(listen); err != nil {
		log.Fatalf("SampleServer failed to serve connections: %v", err)
	}
}

// this struct must implement the methods defined in our Sample Service
type SampleServer struct{}

func (ss *SampleServer) Greet(ctx context.Context, req *pb.Greetee) (resp *pb.GreetResponse, err error) {
	resp = &pb.GreetResponse{Greeting: "Howdy, " + req.FirstName + " " + req.LastName + "!"}
	return
}

func (ss *SampleServer) LifoEcho(stream pb.Sample_LifoEchoServer) error {
	vals := []*pb.Val{}
	intervals := []time.Duration{}
	lastReceived := time.Now()
	for {
		val, err := stream.Recv()
		if err == io.EOF {
			break
		}
		fmt.Printf("Received: %s %d\n", val.GetLabel(), val.GetValue())
		vals = append(vals, val)
		now := time.Now()
		intervals = append(intervals, now.Sub(lastReceived))
		lastReceived = now
	}

	for i := len(vals) - 1; i >= 0; i-- {
		time.Sleep(intervals[i])
		fmt.Printf("Echoing: %s %d\n", vals[i].GetLabel(), vals[i].GetValue())
		stream.Send(vals[i])
	}
	return nil
}
