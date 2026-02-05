# Не знаю, честно говоря, как сюда прокидывать из env...Поставил дефолтные, не ругайте... ಥ_ಥ
DB_URL=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable

migrateUp:
	migrate -path internal/infrastructure/database/postgres/migrations -database "$(DB_URL)" up

migrateDown:
	migrate -path internal/infrastructure/database/postgres/migrations -database "$(DB_URL)" down

swag:
	swag init -g cmd/app/main.go --parseDependency --parseInternal

install-deps:
	go install github.com/vektra/mockery/v2@latest
	go install github.com/swaggo/swag/cmd/swag@latest


mock:
	@echo "Generating mocks using keeptree..."
	mockery --all --recursive --keeptree --dir ./internal --output ./mocks
	@echo "Mocks generated in ./mocks"

mock-clean:
	rm -rf ./mocks