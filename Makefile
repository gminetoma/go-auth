include .env
export

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1

migrate-down-all:
	migrate -path migrations -database "$(DATABASE_URL)" down

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate

gqlgen:
	go tool gqlgen generate
