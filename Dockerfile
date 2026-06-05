FROM golang:tip-alpine3.23 AS builder

COPY / /app

WORKDIR /app/cmd/getweather

RUN go build -o getweather.n .

FROM golang:tip-alpine3.23
COPY --from=builder /app/cmd/getweather/getweather.n /
WORKDIR /
CMD ["./getweather.n"]
