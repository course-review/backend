FROM golang:1.23.6-bookworm

WORKDIR /app
COPY . .

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

RUN sqlc generate

ENTRYPOINT ["go", "run", "server/fiber.go"]