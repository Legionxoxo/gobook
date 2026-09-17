# GoBook

A single-showing cinema seat-booking application built with Go, Redis, REST, and gRPC.

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

The browser uses REST because it works naturally with HTTP and JSON. Before holding a seat, the booking API makes a synchronous unary gRPC call to the internal pricing service. The returned quote is stored with the Redis booking and included in the REST response.

The interface is intentionally limited to one Inception showing, so it focuses on seat holds,
confirmation, expiry, and the internal gRPC pricing call. All seats currently use one mocked
normal-seat price; there are no seat types or dynamic pricing rules.

## Run with Docker

Start the complete containerized application:

```bash
docker compose up --build
```

Then open `http://localhost:8080`. Redis Commander is available at `http://localhost:8081`.

The booking container connects to Redis and pricing through the Compose network using
`redis:6379` and `pricing:9090`; those addresses are provided by `REDIS_ADDR` and
`PRICING_ADDR` environment variables.

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

## Interview explanation

REST is the external API for the browser. gRPC is used for the synchronous internal call from booking to pricing because it provides a strongly typed protobuf contract and generated Go clients. Pricing is kept behind a separate boundary because the same rules could later be shared by web, mobile, and kiosk booking channels.

The RPC is unary because requesting a price is a single request-response operation; streaming would add complexity without helping this use case.

## Test

The project has one Redis-backed stress test. It starts 100,000 concurrent attempts to
hold the same seat and verifies that Redis `SET NX` permits exactly one winner.

Start the Docker stack first, then run:

```bash
go test ./internal/booking -run TestConcurrentBooking_ExactlyOneWins
```
