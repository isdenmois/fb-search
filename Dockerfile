FROM oven/bun:1.3.1-alpine as myjs
FROM golang:1.25.4-alpine AS mygo
FROM alpine as myrun

# install node_modules
FROM myjs AS js-builder
WORKDIR /app
COPY package.json .
COPY bun.lock .
RUN --mount=type=cache,target=/root/.bun/install/cache bun install --frozen-lockfile
COPY . .
RUN bun run web:build

# build go app
FROM mygo AS go-builder
WORKDIR /app
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -ldflags="-s -w" -o main .

FROM myrun
WORKDIR /app
COPY --from=go-builder /app/main .
COPY --from=js-builder /app/public .
CMD ["./main"]
