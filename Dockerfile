# build 7z, jxl
FROM debian:13-slim AS libs

RUN apt-get update && \
    apt-get install -y p7zip-full libjxl-tools && \
    rm -rf /var/lib/apt/lists/*

# install node_modules
FROM oven/bun:1.2.23-slim AS modules
WORKDIR /app
COPY package.json .
COPY bun.lock .

# Mount Bun's cache directory
RUN --mount=type=cache,target=/root/.bun/install/cache bun install --frozen-lockfile

# build the files
FROM oven/bun:1.2.23-slim AS builder
WORKDIR /app
COPY --from=modules /app/node_modules node_modules/
COPY . .
RUN bun run build
RUN bun run web:build

# run the app
FROM libs
WORKDIR /app

ARG PORT=3000
ENV PORT ${PORT}
ENV NODE_ENV production
EXPOSE $PORT
COPY --from=builder /usr/local/bin/bun /usr/bin
COPY --from=builder /app/drizzle drizzle
COPY --from=builder /app/dist dist
COPY --from=builder /app/public public
CMD ["/usr/bin/bun", "dist/index.js"]
