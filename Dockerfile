FROM golang:tip-alpine3.23 AS app-builder
COPY . /app
WORKDIR /app
RUN go build -o getweather.n .

FROM alpine:latest AS cert-builder
RUN apk --no-cache add ca-certificates

FROM scratch
COPY --from=app-builder /app/getweather.n /
COPY --from=cert-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /
CMD ["./getweather.n"]
