package booking

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Legionxoxo/gobook/internal/utils"
)

type handler struct {
	svc    *Service
	pricer PriceQuoter
}

func NewHandler(svc *Service, pricer PriceQuoter) *handler {
	return &handler{svc: svc, pricer: pricer}
}

type holdSeatRequest struct {
	UserID string `json:"user_id"`
}

func (h *handler) HoldSeat(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")
	seatID := r.PathValue("seatID")

	var req holdSeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()

	quote, err := h.pricer.GetQuote(ctx, movieID, seatID)
	if err != nil {
		log.Printf("get price quote: %v", err)
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{
			"error": "pricing service unavailable",
		})
		return
	}

	data := Booking{
		UserID:     req.UserID,
		SeatID:     seatID,
		MovieID:    movieID,
		PricePaise: quote.PricePaise,
		Currency:   quote.Currency,
	}

	session, err := h.svc.Book(data)
	if err != nil {
		log.Println(err)
		return
	}

	type holdResponse struct {
		SessionID  string `json:"session_id"`
		MovieID    string `json:"movie_id"`
		SeatID     string `json:"seat_id"`
		PricePaise int64  `json:"price_paise"`
		Currency   string `json:"currency"`
		ExpiresAt  string `json:"expires_at"`
	}

	utils.WriteJSON(w, http.StatusCreated, holdResponse{
		SeatID:     seatID,
		MovieID:    session.MovieID,
		SessionID:  session.ID,
		PricePaise: session.PricePaise,
		Currency:   session.Currency,
		ExpiresAt:  session.ExpiresAt.Format(time.RFC3339),
	})
}

func (h *handler) ListSeats(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")

	bookings := h.svc.ListBookings(movieID)

	seats := make([]seatInfo, 0, len(bookings))
	for _, b := range bookings {
		seats = append(seats, seatInfo{
			SeatID:    b.SeatID,
			UserID:    b.UserID,
			Booked:    true,
			Confirmed: b.Status == "confirmed",
		})
	}

	utils.WriteJSON(w, http.StatusOK, seats)
}

type seatInfo struct {
	SeatID    string `json:"seat_id"`
	UserID    string `json:"user_id"`
	Booked    bool   `json:"booked"`
	Confirmed bool   `json:"confirmed"`
}

func (h *handler) ConfirmSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("sessionID")

	var req holdSeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return
	}

	if req.UserID == "" {
		return
	}

	session, err := h.svc.ConfirmSeat(r.Context(), sessionID, req.UserID)
	if err != nil {
		return
	}

	utils.WriteJSON(w, http.StatusOK, sessionResponse{
		SessionID: session.ID,
		MovieID:   session.MovieID,
		SeatID:    session.SeatID,
		UserID:    req.UserID,
		Status:    session.Status,
	})
}

type sessionResponse struct {
	SessionID string `json:"session_id"`
	MovieID   string `json:"movie_id"`
	SeatID    string `json:"seat_id"`
	UserID    string `json:"user_id"`
	Status    string `json:"status"`
	ExpiresAt string `json:"expires_at,omitempty"`
}

func (h *handler) ReleaseSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("sessionID")

	var req holdSeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(err)
		return
	}
	if req.UserID == "" {
		return
	}

	err := h.svc.ReleaseSeat(r.Context(), sessionID, req.UserID)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
