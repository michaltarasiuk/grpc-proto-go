package main

import (
	"context"
	"log"
	"net"
	"os"

	pingv1 "github.com/michaltarasiuk/grpc-proto-lab/gen/ping/v1"
	"google.golang.org/grpc"
)

type pingServer struct {
	pingv1.UnimplementedPingServiceServer
}

func (s *pingServer) Ping(ctx context.Context, req *pingv1.PingRequest) (*pingv1.PingResponse, error) {
	msg := req.GetMessage()
	log.Printf("Ping received: %q", msg)
	return &pingv1.PingResponse{Message: "pong: " + msg}, nil
}

func main() {
	addr := envOr("PING_ADDR", ":50051")

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen %s: %v", addr, err)
	}

	s := grpc.NewServer()
	pingv1.RegisterPingServiceServer(s, &pingServer{})

	log.Printf("ping-server listening on %s", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
