# GoBook

A deliberately small, single-showing cinema seat-booking application built with Go, Redis, REST, and gRPC.

## Architecture

```text
Browser
  |
  | REST/JSON :8080
  v
Booking API -------- gRPC :9090 --------> Pricing service
  |                                      GetQuote(movie_id, seat_id)
  v
Redis                                  Fixed normal-seat price: INR 200.00
```

The browser shows one fixed Inception screening. It uses REST because browsers work naturally with HTTP and JSON. Before holding a seat, the Booking API makes a synchronous unary gRPC call to Pricing. The returned quote is saved with the Redis reservation and returned to the browser.

The UI refreshes seat data after a hold, confirmation, release, or expiry; it does not poll. All seats use one mocked normal-seat price, so there are no seat types or dynamic pricing rules.

## REST API

```text
GET    /movies/inception/seats
POST   /movies/inception/seats/{seatID}/hold
PUT    /sessions/{sessionID}/confirm
DELETE /sessions/{sessionID}
```

The hold, confirm, and release routes expect JSON with a `user_id` field.

## Run with Docker

Build and start the complete application:

```bash
docker compose up --build
```

Then open `http://localhost:8080`. Redis Commander is available at `http://localhost:8081`.

The Booking container connects to Redis and Pricing through the Compose network using `redis:6379` and `pricing:9090`. Those addresses are supplied through `REDIS_ADDR` and `PRICING_ADDR` environment variables.

## Run locally

Start Redis:

```bash
docker compose up -d redis
```

Start the pricing gRPC service:

```bash
go run ./cmd/pricing
```

In another terminal, start the booking API and website:

```bash
go run ./cmd
```

Open `http://localhost:8080`.

## gRPC contract

The service exposes one unary RPC:

```text
PricingService.GetQuote(movie_id, seat_id) -> price_paise, currency
```

The contract lives in `api/pricing/v1/pricing.proto`. Generated Go client and server bindings live in `gen/pricing/v1`.

Money is represented as integer paise rather than floating point. For example, INR 200.00 is stored as `20000` paise.

## Project structure

```text
cmd/main.go                    Booking REST API entry point
cmd/pricing/main.go            Pricing gRPC service entry point
internal/booking/              HTTP handlers and Redis reservation store
internal/adapters/pricing/     gRPC client adapter used by Booking
internal/pricing/              Pricing gRPC implementation
api/pricing/v1/                Protobuf contract
gen/pricing/v1/                Generated protobuf and gRPC Go code
static/index.html              Single-showing browser UI
```

## Interview explanation

REST is the external API for the browser. gRPC is used for the synchronous internal call from Booking to Pricing because it provides a strongly typed protobuf contract and generated Go clients. Redis uses an atomic `SET NX` operation with a two-minute TTL, so only one user can hold a seat and abandoned holds expire automatically.

The RPC is unary because requesting a price is a single request-response operation; streaming would add complexity without helping this use case.

## Test

The project has one Redis-backed stress test in `internal/booking/redis_store_test.go`. It starts 100,000 concurrent attempts to hold the same seat and verifies that Redis `SET NX` permits exactly one winner.

Start the Docker stack first, then run:

```bash
go test ./internal/booking -run TestConcurrentBooking_ExactlyOneWins
```
