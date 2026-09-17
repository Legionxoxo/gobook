package booking

import (
	"context"
	"errors"
	"time"
)

var (
	ErrSeatAlreadyBooked = errors.New("seat is already taken")
)

// Booking represents one seat reservation. Its status is "held" until payment
// succeeds, then it becomes "confirmed" and no longer expires.
type Booking struct {
	ID         string
	MovieID    string
	SeatID     string
	UserID     string
	Status     string
	PricePaise int64
	Currency   string
	ExpiresAt  time.Time
}

type BookingStore interface {
	// BookingStore lets the service use Redis in production and a fake store in tests.
	Book(b Booking) (Booking, error)
	ListBookings(movieID string) []Booking

	Confirm(ctx context.Context, sessionID string, userID string) (Booking, error)
	Release(ctx context.Context, sessionID string, userID string) error
}
