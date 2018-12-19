package main

import (
	"fmt"
	"io"
	"log"

	helloService "github.com/wzulfikar/grpc-demo/codegen/go/services/hello"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const addr = "localhost:50000"

func main() {
	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("couldn't connect to grpc endpoint at %s: %v", addr, err)
	}
	defer conn.Close()

	log.Println("connected to grpc endpoint at", addr)

	// Creates a new CustomerClient
	client := helloService.NewHelloServiceClient(conn)

	request := &helloService.HelloRequest{
		SenderName: "Gopher",
	}

	resp, err := client.Greet(context.Background(), request)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.GetGreeting())

	go getStream(client)

	done := make(chan struct{}, 0)
	<-done
}

func getStream(client helloService.HelloServiceClient) {
	log.Println("initializing stream request..")
	request := &helloService.Empty{}

	// calling the streaming API
	stream, err := client.GetStream(context.Background(), request)
	if err != nil {
		log.Fatalf("Error on get customers: %v", err)
	}

	for {
		streamContent, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("something went wrong during streaming: %v", err)
		}
		log.Println("stream counter:", streamContent.GetCounter())
	}
}
