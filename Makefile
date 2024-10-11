build:
	@go build -o bin/gator ./cmd/gator/

run: build
	@./bin/gator

test:
	@go test -v ./...

install-tools:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest


migrate-up:
	@goose --dir ./sql/schema/ postgres "postgres://andrewgaskell:@localhost:5432/gator" up

migrate-down:
	@goose --dir ./sql/schema/ postgres "postgres://andrewgaskell:@localhost:5432/gator" down
