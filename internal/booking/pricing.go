package booking

import "context"

type PriceQuote struct {
	// PricePaise uses an integer to avoid floating-point money errors.
	PricePaise int64
	Currency   string
}

type PriceQuoter interface {
	// PriceQuoter keeps booking independent from gRPC and easy to unit test.
	GetQuote(ctx context.Context, movieID, seatID string) (PriceQuote, error)
}
