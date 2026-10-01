# Multi-stage build for Go backend
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/api cmd/api/main.go

# Production minimal container
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/api /app/api
COPY --from=builder /app/migrations /app/migrations

RUN mkdir -p /app/uploads

EXPOSE 8080

CMD ["/app/api"]
