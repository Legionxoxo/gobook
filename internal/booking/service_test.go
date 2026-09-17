package booking

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/Legionxoxo/gobook/internal/adapters/redis"
	"github.com/google/uuid"
)

func TestConcurrentBooking_ExactlyOneWins(t *testing.T) {
	store := NewRedisStore(redis.NewClient("localhost:6379"))
	svc := NewService(store)

	const numGoroutines = 100_000 // 100k users trying to book a seat at the same time
	movieID := "concurrency-" + uuid.NewString()
	seatID := "A1"

	var (
		successes        atomic.Int64
		expectedFailures atomic.Int64
		unexpectedErrors atomic.Int64
		winnerID         atomic.Pointer[string]
		wg               sync.WaitGroup
	)

	wg.Add(numGoroutines)
	for i := range numGoroutines {
		go func(userNum int) {
			defer wg.Done()
			booking, err := svc.Book(Booking{
				MovieID: movieID,
				SeatID:  seatID,
				UserID:  fmt.Sprintf("user-%d", userNum),
			})
			switch {
			case err == nil:
				successes.Add(1)
				id := booking.ID
				winnerID.Store(&id)
			case errors.Is(err, ErrSeatAlreadyBooked):
				expectedFailures.Add(1)
			default:
				unexpectedErrors.Add(1)
			}
		}(i)
	}
	wg.Wait()

	// The test uses a unique Redis key, then cleans both sides of the session lookup.
	keys := []string{fmt.Sprintf("seat:%s:%s", movieID, seatID)}
	if id := winnerID.Load(); id != nil {
		keys = append(keys, sessionKey(*id))
	}
	defer store.rdb.Del(context.Background(), keys...)

	if got := successes.Load(); got != 1 {
		t.Errorf("expected exactly 1 success, got %d", got)
	}
	if got := expectedFailures.Load(); got != int64(numGoroutines-1) {
		t.Errorf("expected %d seat-taken errors, got %d", numGoroutines-1, got)
	}
	if got := unexpectedErrors.Load(); got != 0 {
		t.Errorf("expected no unexpected errors, got %d", got)
	}
}
