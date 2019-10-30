package main

import (
	"context"
	"fmt"
	"grpc-go/calculator/calculatorpb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Caclulator client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Couldn't connect: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	//fmt.Printf("Created client: %f", c)

	// doUnary(c)

	//doServerStreaming(c)

	doClientStreaming(c)
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while opening stream : %v", err)
	}

	numbers := []int32{3, 5, 9, 56, 15}

	for _, number := range numbers {
		fmt.Printf("Sending number: %v\n", number)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response: %v", err)
	}

	fmt.Printf("The average is: %v\n", res.GetAverage())

}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {

	fmt.Println("Starting to do a PrimeDecomposition  Server Streaming RPC")

	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 23,
	}

	stream, err := c.PrimeNumberDecomposition(context.Background(), req)

	if err != nil {
		log.Fatalf("error while calling Sum RPC %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Something happened %v", err)
		}
		fmt.Println(res.GetPrimeFactor())

	}

}

func doUnary(c calculatorpb.CalculatorServiceClient) {

	fmt.Println("Starting to do a Sum Unary RPC")

	req := &calculatorpb.SumRequest{
		FirstNumber:  3,
		SecondNumber: 20,
	}

	res, err := c.Sum(context.Background(), req)

	if err != nil {
		log.Fatalf("error while calling Sum RPC %v", err)
	}
	log.Printf("Response from Sum: %v", res.SumResult)

}
