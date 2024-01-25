DB_URL := "postgres://postgres:nodirbek@localhost:5432/migration?sslmode=disable"

migrate-up:
	migrate -path migrations -database "$(DB_URL)" -verbose up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" -verbose down

migrate-file:
	migrate create -ext sql -dir migrations/ -seq relationship
