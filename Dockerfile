FROM golang:1.23-alpine AS builder
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o finance-api ./cmd/finance-api

FROM alpine:latest
COPY --from=builder /app/finance-api /finance-api
COPY .env.example .
CMD ["./finance-api"]