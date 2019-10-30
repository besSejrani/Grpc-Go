package main

import (
	"context"
	"fmt"
	"grpc-go/greet/greetpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("hello world")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Couldn't connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	//fmt.Printf("Created client: %f", c)

	// Unary(c)

	//doServerStreaming(c)

	doClientStreaming(c)
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Bes",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Asan",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Alex",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Pietro",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Tina",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling long greet", err)
	}

	//we iterate over our slice and send each message individually
	for _, req := range requests {
		fmt.Printf("sending req: %v\n", req)

		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from LongGreet %v", err)
	}

	fmt.Printf("LongGreat Response: %v\n", res)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Bes",
			LastName:  "Sejrani",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}
}

func doUnary(c greetpb.GreetServiceClient) {

	fmt.Println("Starting to do Unary RPC")

	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Bes",
			LastName:  "Sejrani",
		},
	}

	res, err := c.Greet(context.Background(), req)

	if err != nil {
		log.Fatalf("error while calling Greet RPC %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)

}
