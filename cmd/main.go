package main

import (
	"log"
	"net/http"

	pricingv1 "github.com/Legionxoxo/gobook/gen/pricing/v1"
	pricingadapter "github.com/Legionxoxo/gobook/internal/adapters/pricing"
	"github.com/Legionxoxo/gobook/internal/adapters/redis"
	"github.com/Legionxoxo/gobook/internal/booking"
	"github.com/Legionxoxo/gobook/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /movies", listMovies)

	mux.Handle("GET /", http.FileServer(http.Dir("static")))

	store := booking.NewRedisStore(redis.NewClient("localhost:6379"))
	svc := booking.NewService(store)
	pricingConn, err := grpc.NewClient("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer pricingConn.Close()

	priceClient := pricingadapter.NewClient(pricingv1.NewPricingServiceClient(pricingConn))
	bookingHandler := booking.NewHandler(svc, priceClient)

	mux.HandleFunc("GET /movies/{movieID}/seats", bookingHandler.ListSeats)
	mux.HandleFunc("POST /movies/{movieID}/seats/{seatID}/hold", bookingHandler.HoldSeat)

	mux.HandleFunc("PUT /sessions/{sessionID}/confirm", bookingHandler.ConfirmSession)
	mux.HandleFunc("DELETE /sessions/{sessionID}", bookingHandler.ReleaseSession)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

var movies = []movieResponse{
	{ID: "inception", Title: "Inception", Rows: 5, SeatsPerRow: 8},
	{ID: "dune", Title: "Dune: Part Two", Rows: 4, SeatsPerRow: 6},
}

func listMovies(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, movies)
}

type movieResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Rows        int    `json:"rows"`
	SeatsPerRow int    `json:"seats_per_row"`
}
