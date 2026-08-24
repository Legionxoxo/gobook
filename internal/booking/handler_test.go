package booking

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type handlerTestStore struct {
	booked Booking
}

func (s *handlerTestStore) Book(b Booking) (Booking, error) {
	s.booked = b
	b.ID = "session-1"
	b.Status = "held"
	b.ExpiresAt = time.Now().Add(2 * time.Minute)
	return b, nil
}

func (s *handlerTestStore) ListBookings(string) []Booking {
	return nil
}

func (s *handlerTestStore) Confirm(context.Context, string, string) (Booking, error) {
	return Booking{}, nil
}

func (s *handlerTestStore) Release(context.Context, string, string) error {
	return nil
}

type handlerTestPricer struct{}

func (handlerTestPricer) GetQuote(context.Context, string, string) (PriceQuote, error) {
	return PriceQuote{PricePaise: 20_000, Currency: "INR"}, nil
}

func TestHoldSeatIncludesPrice(t *testing.T) {
	store := &handlerTestStore{}
	handler := NewHandler(NewService(store), handlerTestPricer{})
	mux := http.NewServeMux()
	mux.HandleFunc("POST /movies/{movieID}/seats/{seatID}/hold", handler.HoldSeat)

	req := httptest.NewRequest(http.MethodPost, "/movies/inception/seats/A1/hold", bytes.NewBufferString(`{"user_id":"user-1"}`))
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusCreated)
	}

	var response struct {
		MovieID    string `json:"movie_id"`
		PricePaise int64  `json:"price_paise"`
		Currency   string `json:"currency"`
	}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.MovieID != "inception" || response.PricePaise != 20_000 || response.Currency != "INR" {
		t.Fatalf("response = %+v", response)
	}
	if store.booked.PricePaise != 20_000 || store.booked.Currency != "INR" {
		t.Fatalf("stored booking = %+v", store.booked)
	}
}
