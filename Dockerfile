FROM golang:1.26-alpine

WORKDIR /app
COPY . .

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

RUN sqlc generate

WORKDIR /app/server

RUN go mod download

ENTRYPOINT ["go", "run", "."]