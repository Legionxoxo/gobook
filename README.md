# GoBook

A small cinema seat-booking application built with Go, Redis, REST, and gRPC.

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

All seats currently use one mocked normal-seat price. There are no seat types or dynamic pricing rules.

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

## Tests

The pricing tests include validation, generated-client mapping, and a real gRPC client/server call over an in-memory transport:

```bash
go test ./internal/pricing ./internal/adapters/pricing
go test ./internal/booking -run TestHoldSeatIncludesPrice
```

The existing concurrent Redis test requires Redis to be running.
