package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	helloService "github.com/wzulfikar/grpc-demo/codegen/go/services/hello"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// `server` implements `helloService` proto
type server struct{}

func logRequest(methodName string, ctx context.Context, params interface{}) {
	p, _ := peer.FromContext(ctx)
	meta, _ := metadata.FromIncomingContext(ctx)

	log.Printf("server.%s: ip_addr=%s user_agent=%s params=%v", methodName, p.Addr, meta["user-agent"], params)
}

func (s *server) Greet(ctx context.Context, in *helloService.HelloRequest) (*helloService.HelloResponse, error) {
	logRequest("Greet", ctx, in)

	greeting := fmt.Sprintf("Hello %s, welcome to grpc!", in.SenderName)
	return &helloService.HelloResponse{Greeting: greeting}, nil
}

func (s *server) GetStream(empty *helloService.Empty, stream helloService.HelloService_GetStreamServer) error {
	logRequest("GetStream", stream.Context(), nil)

	counter := &helloService.StreamResponse{Counter: 0}
	for {
		if err := stream.Send(counter); err != nil {
			return err
		}
		time.Sleep(time.Second * 1)
		counter.Counter++
	}
	return nil
}

var addr = fmt.Sprintf(":%s", os.Getenv("PORT"))

func main() {
	tcpListen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}

	log.Println("grpc listening at", addr)

	// create new gRPC server
	s := grpc.NewServer()
	helloService.RegisterHelloServiceServer(s, &server{})
	s.Serve(tcpListen)
}
