# Stage 1: Build the frontend
FROM node:22-alpine AS frontend-builder
WORKDIR /web
COPY web/package*.json ./
RUN npm install
COPY web/ ./
RUN npm run build

# Stage 2: Build the Go backend
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o brokeli cmd/main.go

# Stage 3: Final image
FROM alpine:latest
WORKDIR /app

# Copy the Go binary
COPY --from=builder /app/brokeli .

# Copy the built frontend static files
# The Go backend is configured to serve them from web/dist
COPY --from=frontend-builder /web/dist ./web/dist

# Copy other required assets
COPY internal/features/import_transactions/transactions.csv internal/features/import_transactions/transactions.csv

EXPOSE 8080

CMD ["./brokeli"]
