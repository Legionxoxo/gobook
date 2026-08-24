package pricing

import (
	"context"
	"net"
	"testing"

	pricingv1 "github.com/Legionxoxo/gobook/gen/pricing/v1"
	pricingservice "github.com/Legionxoxo/gobook/internal/pricing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func TestClientCallsPricingServer(t *testing.T) {
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	pricingv1.RegisterPricingServiceServer(server, pricingservice.NewServer())
	go func() {
		if err := server.Serve(listener); err != nil {
			t.Errorf("Serve() error = %v", err)
		}
	}()
	t.Cleanup(server.Stop)

	conn, err := grpc.NewClient(
		"passthrough:///pricing",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	quote, err := NewClient(pricingv1.NewPricingServiceClient(conn)).GetQuote(
		context.Background(),
		"inception",
		"A1",
	)
	if err != nil {
		t.Fatalf("GetQuote() error = %v", err)
	}
	if quote.PricePaise != 20_000 || quote.Currency != "INR" {
		t.Fatalf("GetQuote() = %+v", quote)
	}
}
