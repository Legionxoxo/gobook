package pricing

import (
	"context"
	"strings"

	pricingv1 "github.com/Legionxoxo/gobook/gen/pricing/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const normalSeatPricePaise int64 = 20_000

type Server struct {
	// Embedding preserves forward compatibility when RPCs are added to the proto.
	pricingv1.UnimplementedPricingServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) GetQuote(_ context.Context, req *pricingv1.GetQuoteRequest) (*pricingv1.GetQuoteResponse, error) {
	// Validate input at the service boundary before applying pricing rules.
	if req == nil || strings.TrimSpace(req.GetMovieId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "movie_id is required")
	}
	if strings.TrimSpace(req.GetSeatId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "seat_id is required")
	}

	// This is deliberately fixed for now; future rules can use movie/seat inputs.
	return &pricingv1.GetQuoteResponse{
		PricePaise: normalSeatPricePaise,
		Currency:   "INR",
	}, nil
}
