FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app/main.go

#=========================
FROM alpine:3.19

RUN adduser -D appuser
USER appuser

WORKDIR /app

COPY --from=builder /app/main .

COPY --from=builder /app/internal/infrastructure/database/postgres/migrations ./internal/infrastructure/database/postgres/migrations


EXPOSE 8080
CMD ["./main"]