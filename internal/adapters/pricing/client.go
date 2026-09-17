package pricing

import (
	"context"

	pricingv1 "github.com/Legionxoxo/gobook/gen/pricing/v1"
	"github.com/Legionxoxo/gobook/internal/booking"
)

type Client struct {
	// client is generated from api/pricing/v1/pricing.proto.
	client pricingv1.PricingServiceClient
}

func NewClient(client pricingv1.PricingServiceClient) *Client {
	return &Client{client: client}
}

func (c *Client) GetQuote(ctx context.Context, movieID, seatID string) (booking.PriceQuote, error) {
	// Translate booking-domain inputs into the protobuf request expected by gRPC.
	quote, err := c.client.GetQuote(ctx, &pricingv1.GetQuoteRequest{
		MovieId: movieID,
		SeatId:  seatID,
	})
	if err != nil {
		return booking.PriceQuote{}, err
	}

	// Translate the protobuf response back into a booking-domain value.
	return booking.PriceQuote{
		PricePaise: quote.GetPricePaise(),
		Currency:   quote.GetCurrency(),
	}, nil
}
