# FB Search

A book search application for the Flibusta library, enabling search and download of FB2 books via a modern web interface. The project uses a hybrid Go/TypeScript architecture with PostgreSQL for data storage.

## Architecture

**Dual-Stack Architecture:**

1. **Backend (Go)**: REST API server using Gin framework
   - Entry point: `main.go`
   - Dependency injection via `sarulabs/di`
   - PostgreSQL database via `pgx/v5`
   - HTTP controllers in `views/controllers/`

2. **Frontend (Vue 3)**: SPA using Vite + UnoCSS
   - Entry point: `web/app/main.ts`
   - Vue Router for client-side routing
   - HTTP API client via `wretch`

### Data Flow

```
User Browser (Vue SPA)
        ↓
   Vite Dev Server / HTTP API
        ↓
  PostgreSQL Database (books index)
        ↓
  Flat Files (ZIP archives in files/)
```

## Tech Stack

- **Backend**: Go 1.25+ (Gin framework, pgx/v5, sarulabs/di)
- **Frontend**: Vue 3 + TypeScript + Vite 7.3.1 + UnoCSS 66.6.2
- **Database**: PostgreSQL
- **Build Tools**: Bun, Go modules
- **Testing**: Playwright (E2E), Vitest (unit), testcontainers-go (integration)

## Project Structure

```
.
├── app/              # Business logic (INP parser)
├── domain/           # Domain entities (Book, ParseProgress)
├── infra/            # Infrastructure layer
│   ├── db/          # Database connection, migrations
│   └── repositories/ # Data access layer
├── migrations/       # PostgreSQL schema migrations
├── parser/           # Book parsing logic
├── views/            # HTTP layer
│   ├── controllers/  # API endpoints
│   ├── di.go        # Dependency injection setup
│   └── views.go     # HTTP server initialization
├── tests/            # Go integration tests
│   ├── integration/  # Controller and repository tests
│   ├── testhelpers/  # Test infrastructure
│   ├── fixtures/     # Test data
│   └── mocks/        # Mock implementations
├── scripts/          # Build and migration scripts
├── web/              # Vue frontend
│   ├── app/         # Vue entry point, App.vue
│   ├── pages/       # Page components
│   ├── entities/    # API types and functions
│   └── shared/      # Shared utilities
├── playwright/       # E2E test suite
├── files/            # Book archives (ZIP, INPX)
└── public/           # Built frontend assets
```

## Getting Started

### Prerequisites

- Go 1.25+
- PostgreSQL
- Bun

### Setup

1. Clone the repository
2. Configure environment variables:

```bash
cp .env.example .env
# Edit .env with your database credentials
```

3. Run database migrations:

```bash
# Frontend
bun run schema

# Backend
go run main.go
```

### Development

```bash
# Frontend (with hot reload)
bun run web

# Backend
go run main.go

# Or via just (if installed)
just run
just build
```

### Importing Books

Place your Flibusta files in the `files/` directory:

- `flibusta_fb2_local.inpx` - Index file
- `*.zip` - Book archives

Parse the books into the database:

```bash
# From frontend
bun run parse

# Or via API
curl -X POST http://localhost:8080/api/parse
```

## Testing

### Go Integration Tests

```bash
go test ./tests/integration/... -v
```

### Unit Tests (Vitest)

```bash
bun run test
```

### E2E Tests (Playwright)

```bash
bun run playwright
```

## API Endpoints

- `GET /api/search?q=<query>` - Search books
- `GET /dl/:file/:path` - Download book file
- `POST /api/parse` - Rebuild database from INPX files

## Docker Deployment

```bash
docker-compose up -d  # Starts PostgreSQL
```

## License

MIT
