# Building
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/dyndns ./cmd/dyndns

# Running
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/dyndns /app/dyndns

RUN adduser -D -H -h /app dyndns && \
    chown -R dyndns:dyndns /app

USER dyndns

ENTRYPOINT ["/app/dyndns"]
