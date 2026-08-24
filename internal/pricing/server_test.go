package pricing

import (
	"context"
	"testing"

	pricingv1 "github.com/Legionxoxo/gobook/gen/pricing/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetQuoteReturnsNormalSeatPrice(t *testing.T) {
	quote, err := NewServer().GetQuote(context.Background(), &pricingv1.GetQuoteRequest{
		MovieId: "inception",
		SeatId:  "A1",
	})
	if err != nil {
		t.Fatalf("GetQuote() error = %v", err)
	}
	if quote.GetPricePaise() != 20_000 {
		t.Fatalf("GetQuote() price = %d, want 20000", quote.GetPricePaise())
	}
	if quote.GetCurrency() != "INR" {
		t.Fatalf("GetQuote() currency = %q, want INR", quote.GetCurrency())
	}
}

func TestGetQuoteValidatesRequest(t *testing.T) {
	tests := []struct {
		name string
		req  *pricingv1.GetQuoteRequest
	}{
		{name: "missing request"},
		{name: "missing movie", req: &pricingv1.GetQuoteRequest{SeatId: "A1"}},
		{name: "missing seat", req: &pricingv1.GetQuoteRequest{MovieId: "inception"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewServer().GetQuote(context.Background(), tt.req)
			if status.Code(err) != codes.InvalidArgument {
				t.Fatalf("GetQuote() code = %v, want InvalidArgument", status.Code(err))
			}
		})
	}
}
