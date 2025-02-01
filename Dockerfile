FROM golang:1.23-alpine AS builder

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o finance-api ./cmd/finance-api

FROM alpine:latest

RUN apk add --no-cache bash postgresql-client

COPY --from=builder /app/finance-api /finance-api
COPY --from=builder /go/bin/goose /usr/local/bin/goose

COPY migrations /migrations
COPY .env.example .


CMD ["sh", "-c", " until pg_isready -h $DB_HOST -p $DB_PORT -U $DB_USER; do echo 'Waiting for database...'; sleep 2; done; goose -dir /migrations postgres \"$DB_CONN\" up && ./finance-api"]
