FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o brokeli cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/brokeli .
COPY internal/features/import_transactions/transactions.csv internal/features/import_transactions/transactions.csv

EXPOSE 8080

CMD ["./brokeli"]
