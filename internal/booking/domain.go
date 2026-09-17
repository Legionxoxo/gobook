package booking

import (
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
