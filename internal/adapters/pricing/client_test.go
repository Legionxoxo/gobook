package pricing

import (
	"context"
	"testing"

	pricingv1 "github.com/Legionxoxo/gobook/gen/pricing/v1"
	"google.golang.org/grpc"
)

type fakePricingClient struct {
	req *pricingv1.GetQuoteRequest
}

func (f *fakePricingClient) GetQuote(_ context.Context, req *pricingv1.GetQuoteRequest, _ ...grpc.CallOption) (*pricingv1.GetQuoteResponse, error) {
	f.req = req
	return &pricingv1.GetQuoteResponse{PricePaise: 20_000, Currency: "INR"}, nil
}

func TestClientGetsQuote(t *testing.T) {
	grpcClient := &fakePricingClient{}
	quote, err := NewClient(grpcClient).GetQuote(context.Background(), "inception", "A1")
	if err != nil {
		t.Fatalf("GetQuote() error = %v", err)
	}
	if grpcClient.req.GetMovieId() != "inception" || grpcClient.req.GetSeatId() != "A1" {
		t.Fatalf("GetQuote() request = %+v", grpcClient.req)
	}
	if quote.PricePaise != 20_000 || quote.Currency != "INR" {
		t.Fatalf("GetQuote() = %+v", quote)
	}
}
