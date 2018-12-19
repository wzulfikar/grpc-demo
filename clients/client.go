package main

import (
	"fmt"
	"io"
	"log"
	"os"

	helloService "github.com/wzulfikar/grpc-demo/codegen/go/services/hello"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var endpoint = os.Getenv("ENDPOINT")

func main() {
	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("couldn't connect to grpc endpoint at %s: %v", endpoint, err)
	}
	defer conn.Close()

	log.Println("connected to grpc endpoint at", endpoint)

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
