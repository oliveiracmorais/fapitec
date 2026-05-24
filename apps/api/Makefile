.PHONY: run test test/coverage build clean vet sqlc

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
