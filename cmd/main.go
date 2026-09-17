package main

import (
	"log"
	"net/http"
	"os"

	pricingv1 "github.com/Legionxoxo/gobook/gen/pricing/v1"
	pricingadapter "github.com/Legionxoxo/gobook/internal/adapters/pricing"
	"github.com/Legionxoxo/gobook/internal/adapters/redis"
	"github.com/Legionxoxo/gobook/internal/booking"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// --- HTTP routes and static website ---
	mux := http.NewServeMux()

	mux.Handle("GET /", http.FileServer(http.Dir("static")))

	// --- Infrastructure dependencies ---
	// Docker supplies service names; local development uses the localhost defaults.
	store := booking.NewRedisStore(redis.NewClient(envOrDefault("REDIS_ADDR", "localhost:6379")))
	svc := booking.NewService(store)
	pricingConn, err := grpc.NewClient(envOrDefault("PRICING_ADDR", "localhost:9090"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer pricingConn.Close()

	// The adapter keeps generated gRPC/protobuf types out of the booking package.
	priceClient := pricingadapter.NewClient(pricingv1.NewPricingServiceClient(pricingConn))
	bookingHandler := booking.NewHandler(svc, priceClient)

	mux.HandleFunc("GET /movies/{movieID}/seats", bookingHandler.ListSeats)
	mux.HandleFunc("POST /movies/{movieID}/seats/{seatID}/hold", bookingHandler.HoldSeat)

	mux.HandleFunc("PUT /sessions/{sessionID}/confirm", bookingHandler.ConfirmSession)
	mux.HandleFunc("DELETE /sessions/{sessionID}", bookingHandler.ReleaseSession)

	// --- Start the public REST API ---
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func envOrDefault(key, fallback string) string {
	// Environment variables make the same binary work locally and in Docker.
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
