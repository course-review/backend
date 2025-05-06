To generate the sqlc files:  
`sqlc generate`

To run the server:  
`docker compose up -d`  
`go run server/fiber.go`

docker:
have postgres, sqlc and go
set up db with schema.sql
sqlc generate
go run server/fiber.go

You also need a logs directory under server/