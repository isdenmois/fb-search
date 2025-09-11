# build 7z, jxl
FROM debian:13-slim AS libs

RUN apt-get update && \
    apt-get install -y p7zip-full libjxl-tools && \
    rm -rf /var/lib/apt/lists/*

COPY --from=oven/bun:slim /usr/local/bin/bun /usr/local/bin

# install node_modules
FROM oven/bun:slim AS modules
WORKDIR /app
COPY package.json .
COPY bun.lockb .
RUN bun install

# build the files
FROM oven/bun:slim AS builder
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
COPY package.json .
COPY tsconfig.json .
COPY --from=builder /app/dist dist
COPY --from=builder /app/public public
COPY --from=builder /app/drizzle drizzle
CMD ["bun", "start"]
