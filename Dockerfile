FROM golang:tip-alpine3.23 AS builder

COPY / /app

WORKDIR /app

RUN go build -o getweather.n .

FROM golang:tip-alpine3.23
COPY --from=builder /app/getweather.n /
WORKDIR /
CMD ["./getweather.n"]
