.PHONY: setup assets migrate seed run test build
setup:
	go mod download
	npm install
	npm run build
assets:
	npm run build
migrate:
	go run ./cmd/migrate
seed:
	go run ./cmd/seed
run:
	go run ./cmd/server
test:
	go test ./...
build:
	go build -o bin/lms ./cmd/server
