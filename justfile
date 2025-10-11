build:
    go build  -ldflags="-s -w" -o dist/fb-search parser/*.go

run:
    go run parser/*.go
