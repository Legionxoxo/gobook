FROM golang:1.25.3-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/booking ./cmd \
    && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/pricing ./cmd/pricing

FROM alpine:3.22 AS booking

WORKDIR /app
COPY --from=build /out/booking ./booking
COPY static ./static

EXPOSE 8080
ENTRYPOINT ["./booking"]

FROM alpine:3.22 AS pricing

COPY --from=build /out/pricing /pricing

EXPOSE 9090
ENTRYPOINT ["/pricing"]
