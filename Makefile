# Не знаю, честно говоря, как сюда прокидывать из env...Поставил дефолтные, не ругайте... ಥ_ಥ
DB_URL=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable

migrateUp:
	migrate -path internal/infrastructure/database/postgres/migrations -database "$(DB_URL)" up

migrateDown:
	migrate -path internal/infrastructure/database/postgres/migrations -database "$(DB_URL)" down

swag:
	swag init -g cmd/app/main.go --parseDependency --parseInternal