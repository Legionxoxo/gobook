package pricing

import (
	"context"

	pricingv1 "github.com/Legionxoxo/gobook/gen/pricing/v1"
	"github.com/Legionxoxo/gobook/internal/booking"
)

type Client struct {
	client pricingv1.PricingServiceClient
}

func NewClient(client pricingv1.PricingServiceClient) *Client {
	return &Client{client: client}
}

func (c *Client) GetQuote(ctx context.Context, movieID, seatID string) (booking.PriceQuote, error) {
	quote, err := c.client.GetQuote(ctx, &pricingv1.GetQuoteRequest{
		MovieId: movieID,
		SeatId:  seatID,
	})
	if err != nil {
		return booking.PriceQuote{}, err
	}

	return booking.PriceQuote{
		PricePaise: quote.GetPricePaise(),
		Currency:   quote.GetCurrency(),
	}, nil
}
