package booking

import "context"

type Service struct {
	// store owns persistence and seat-level concurrency rules.
	store BookingStore
}

func NewService(store BookingStore) *Service {
	return &Service{store}
}

func (s *Service) Book(b Booking) (Booking, error) {
	// Delegate the hold operation to the configured persistence implementation.
	return s.store.Book(b)
}

func (s *Service) ListBookings(movieID string) []Booking {
	return s.store.ListBookings(movieID)
}

func (s *Service) ConfirmSeat(ctx context.Context, sessionID string, userID string) (Booking, error) {
	// Confirmation makes a temporary hold permanent.
	return s.store.Confirm(ctx, sessionID, userID)
}

func (s *Service) ReleaseSeat(ctx context.Context, sessionID string, userID string) error {
	// Releasing frees a seat before its hold TTL expires.
	return s.store.Release(ctx, sessionID, userID)
}
