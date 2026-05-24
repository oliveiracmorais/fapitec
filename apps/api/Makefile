.PHONY: run test test/coverage build clean vet sqlc db/up db/down db/migrate db/psql

run:
	go run ./cmd/api

test:
	go test ./... -v -count=1

test/coverage:
	go test ./... -coverprofile=coverage.out -count=1
	go tool cover -html=coverage.out -o coverage.html
	@echo "Cobertura gerada em coverage.html"

build:
	go build -o bin/api ./cmd/api

clean:
	rm -rf bin/ coverage.out coverage.html

vet:
	go vet ./...

sqlc:
	@command -v sqlc >/dev/null 2>&1 && sqlc generate || echo "sqlc nao instalado. Instale em: https://sqlc.dev"

db/up:
	@echo "Iniciando PostgreSQL..."
	docker compose up -d
	@echo "Aguardando banco ficar pronto..."
	@sleep 3
	@echo "Banco pronto em localhost:5432"

db/down:
	docker compose down

db/psql:
	@command -v psql >/dev/null 2>&1 && PGPASSWORD=fapitec123 psql -h localhost -U fapitec -d fapitec || echo "psql nao instalado. Use: docker exec -it fapitec-db psql -U fapitec"
