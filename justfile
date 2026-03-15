set dotenv-load := true

[parallel]
run: run-client run-go

run-go:
    cd server && go run main.go

run-client:
    bun run dev

run-db:
    docker compose up

build:
    cd server && go build -ldflags="-s -w" -o ../dist/fb-search .

[parallel]
test: test-go test-vitest test-pw

test-go:
    cd server && go test ./...

test-vitest:
    bun run test

test-pw:
    bun run playwright

