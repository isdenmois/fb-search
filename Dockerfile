# install node_modules
FROM oven/bun:latest AS modules
WORKDIR /app
COPY package.json .
COPY bun.lockb .
RUN bun install

# build the files
FROM oven/bun:latest as builder
WORKDIR /app
COPY --from=modules /app/node_modules node_modules/
COPY . .
RUN bun run build
RUN bun run web:build

# run the app
FROM oven/bun:latest
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
