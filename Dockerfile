FROM golang:1.12 AS builder

RUN mkdir /app

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o=authproxy cmd/authproxy/main.go

FROM alpine:3.9

RUN mkdir /app

WORKDIR /app

COPY --from=builder /app/authproxy /app/authproxy

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["/app/authproxy"]