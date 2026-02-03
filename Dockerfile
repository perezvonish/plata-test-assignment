FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

RUN /go/bin/swag init -g cmd/app/main.go --parseDependency --parseInternal

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app/main.go

#=========================
FROM alpine:3.19

RUN adduser -D appuser
WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/internal/infrastructure/database/postgres/migrations ./internal/infrastructure/database/postgres/migrations
COPY --from=builder /app/docs ./docs

RUN chown -R appuser:appuser /app
USER appuser

EXPOSE 8080
CMD ["./main"]