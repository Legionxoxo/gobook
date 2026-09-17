package main

import (
	"log"
	"net"

	pricingv1 "github.com/Legionxoxo/gobook/gen/pricing/v1"
	"github.com/Legionxoxo/gobook/internal/pricing"
	"google.golang.org/grpc"
)

func main() {
	// --- Start the internal pricing gRPC service ---
	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatal(err)
	}

	// Register the handwritten pricing implementation against the generated contract.
	server := grpc.NewServer()
	pricingv1.RegisterPricingServiceServer(server, pricing.NewServer())

	log.Printf("pricing gRPC server listening on %s", listener.Addr())
	if err := server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
