package booking

import "context"

type PriceQuote struct {
	PricePaise int64
	Currency   string
}

type PriceQuoter interface {
	GetQuote(ctx context.Context, movieID, seatID string) (PriceQuote, error)
}
