FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o brokeli cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/brokeli .
# Copy migrations preserving the path structure expected by the app
COPY --from=builder /app/pkg/event_store/postgres/migrations ./pkg/event_store/postgres/migrations

EXPOSE 8080

CMD ["./brokeli"]
