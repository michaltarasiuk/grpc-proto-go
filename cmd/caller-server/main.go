package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	pingv1 "github.com/michaltarasiuk/grpc-proto-lab/gen/ping/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	httpAddr := envOr("CALLER_HTTP_ADDR", ":8080")
	pingAddr := envOr("PING_ADDR", "localhost:50051")

	conn, err := grpc.NewClient(
		pingAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("dial ping-server at %s: %v", pingAddr, err)
	}
	defer conn.Close()

	client := pingv1.NewPingServiceClient(conn)

	mux := http.NewServeMux()
	mux.HandleFunc("/call", func(w http.ResponseWriter, r *http.Request) {
		msg := r.URL.Query().Get("msg")
		if msg == "" {
			msg = "hello"
		}

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		log.Printf("calling Ping(%q) on %s", msg, pingAddr)
		resp, err := client.Ping(ctx, &pingv1.PingRequest{Message: msg})
		if err != nil {
			http.Error(w, fmt.Sprintf("ping failed: %v", err), http.StatusBadGateway)
			return
		}

		fmt.Fprintf(w, "request=%s\nresponse=%s\n", msg, resp.GetMessage())
	})

	srv := &http.Server{Addr: httpAddr, Handler: mux}

	go func() {
		log.Printf("caller-server HTTP on %s → gRPC %s", httpAddr, pingAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http serve: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Print("caller-server shutting down")
	_ = srv.Shutdown(shutdownCtx)
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
