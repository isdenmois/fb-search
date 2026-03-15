# Repository Guidelines

## Project Overview

**FB Search** is a book search application for the Flibusta library, enabling search and download of FB2 books via a modern web interface. The project uses a hybrid Go/TypeScript architecture with PostgreSQL for data storage.

## Architecture & Data Flow

### Dual-Stack Architecture

The project is a **monorepo** with two distinct stacks:

1. **Backend (Go)**: REST API server using Gin framework
   - Entry point: `server/main.go`
   - Dependency injection via `sarulabs/di`
   - PostgreSQL database via `pgx/v5`
   - HTTP controllers in `server/views/controllers/`

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

**Key components:**

- **Domain layer** (`server/domain/`): Pure Go structs (`Book`, `ParseProgress`)
- **App layer** (`server/app/`): Business logic (INP parser)
- **Infrastructure** (`server/infra/`): Database connections, repositories
- **Views layer** (`server/views/`): HTTP controllers, DI configuration

### Dependency Injection Pattern

Backend uses constructor injection via `sarulabs/di`:

```go
// server/views/di.go
func CreateDi() (di.Container, error) {
    builder, _ := di.NewEnhancedBuilder()
    builder.Add(DbDef)
    builder.Add(BooksRepositoryDef)
    builder.Add(ControllersDef)
    builder.Add(HttpServerDef)
    return builder.Build()
}
```

Controllers receive dependencies via constructor:

```go
type BookController struct {
    booksRepository *repositories.BooksRepository
}

func NewBookController(booksRepository *repositories.BooksRepository) *BookController {
    return &BookController{booksRepository: booksRepository}
}
```

## Key Directories

```
.
├── server/              # Go backend
│   ├── main.go         # Entry point
│   ├── app/            # Business logic (INP parser for Flibusta format)
│   ├── domain/         # Domain entities (Book, ParseProgress)
│   ├── infra/         # Infrastructure layer
│   │   ├── db/        # Database connection, migrations
│   │   └── repositories/ # Data access layer (BooksRepository)
│   ├── migrations/    # PostgreSQL schema migrations
│   ├── parser/        # Book parsing logic
│   ├── shared/        # Shared utilities
│   ├── views/         # HTTP layer
│   │   ├── controllers/ # API endpoints (search, download, parse)
│   │   ├── di.go     # Dependency injection setup
│   │   └── views.go  # HTTP server initialization
│   └── tests/         # Go integration tests
│       ├── integration/ # Controller and repository tests
│       ├── testhelpers/ # Test infrastructure (testcontainers)
│       ├── fixtures/  # Test data (sample books)
│       └── mocks/     # Mock implementations
├── scripts/          # Build and migration scripts (TypeScript)
├── web/              # Vue frontend
│   ├── app/         # Vue entry point, App.vue
│   ├── pages/       # Page components
│   ├── entities/    # API types and functions
│   └── shared/      # Shared utilities
├── playwright/       # E2E test suite
│   ├── tests/       # Test specifications
│   └── pages/       # Page object patterns
├── files/            # Book archives (ZIP, INPX)
└── public/           # Built frontend assets
```

## Development Commands

### Backend (Go)

```bash
# Run development server
cd server && go run main.go

# Build production binary
cd server && go build -ldflags="-s -w" -o ../dist/fb-search .

# Run via just (if installed)
just run
just build
```

### Frontend (TypeScript/Bun)

```bash
# Start dev server with hot reload
bun run web

# Build production assets
bun run web:build

# Run unit tests (Vitest)
bun run test

# Run E2E tests (Playwright)
bun run playwright
```

### Import Books

```bash
# Parse INPX index files (from frontend)
bun run parse

# Or via API
curl -X POST http://localhost:8080/api/parse
```

## Code Conventions & Patterns

### Go Backend

**Naming:**

- Exported types/functions start with uppercase (`Book`, `NewBookController`)
- Private helpers use lowercase (`filterASCII`, `search`)
- Controller methods use lowercase with receiver (`(self BookController) search`)

**Error Handling:**

```go
err := someOperation()
if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
}
```

**HTTP Responses:**

```go
c.JSON(http.StatusOK, books)           // Success
c.JSON(http.StatusBadRequest, ...)    // Client error
c.JSON(http.StatusInternalServerError, ...)  // Server error
```

**Controller Pattern:**

```go
type Controller interface {
    Bind(*gin.Engine) error
}

type BookController struct {
    booksRepository *repositories.BooksRepository
}

func (c BookController) Bind(r *gin.Engine) error {
    r.GET("/api/search", c.search)
    r.GET("/dl/:file/:path", c.downloadFile)
    return nil
}
```

### TypeScript Frontend

**Type Safety:**

- Strict TypeScript with `bundler` mode and `noImplicitAny`
- Types defined in `web/entities/`
- API responses fully typed

**HTTP Client:**

```typescript
// Using wretch for type-safe requests
import { wretch } from "wretch";

const api = wretch(BASE_URL).headers({ "Content-Type": "application/json" });
```

**Component Structure:**

- Pages in `web/pages/`
- Shared utilities in `web/shared/`
- UnoCSS utility classes for styling

## Important Files

### Backend Entry Points

- `server/main.go` - Application bootstrap, DI setup
- `server/views/di.go` - Dependency injection configuration
- `server/views/views.go` - HTTP server initialization

### Frontend Entry Points

- `web/app/main.ts` - Vue app bootstrap
- `web/app/App.vue` - Root component
- `vite.config.ts` - Vite + UnoCSS configuration

### Configuration

- `server/go.mod` - Go dependencies (Gin, pgx, di, migrate)
- `package.json` - TypeScript/Bun dependencies (Vue, Vite, Vitest, Playwright)
- `tsconfig.json` - TypeScript compiler options
- `tsconfig.web.json` - Web-specific config with path aliases
- `.env` - Environment variables (`DATABASE_URL`)
- `playwright.config.ts` - E2E test configuration

### Data Layer

- `server/infra/db/db.go` - PostgreSQL connection, migration runner
- `server/infra/repositories/books_repository.go` - Books data access
- `server/domain/book.go` - Book entity definition
- `server/migrations/` - SQL schema migrations

## Runtime/Tooling Preferences

### Required Runtimes

- **Go**: 1.25+ (backend)
- **Bun**: Latest (TypeScript tooling, scripts)
- **PostgreSQL**: Any recent version

### Package Manager

- **Bun** for TypeScript/JavaScript tooling
- **Go modules** for backend dependencies

### Build Tools

- **Vite 7.3.1** - Frontend bundler
- **UnoCSS 66.6.2** - Atomic CSS engine
- **Playwright** - E2E testing
- **Vitest** - Unit testing

### Code Quality

- **Biome** - Lint and Formatter
- **just** (optional) - CLI command runner

## Testing & QA

### Go Integration Tests

Located in `tests/` directory:

- `tests/integration/` - HTTP controller and repository integration tests
- `tests/testhelpers/` - Test infrastructure (testcontainers database)
- `tests/fixtures/` - Test data (sample books)
- `tests/mocks/` - Mock implementations for testing

```bash
go test ./tests/integration/... -v
```

**Features:**

- Uses **testcontainers-go** to spin up PostgreSQL containers automatically
- Tests full HTTP request-response cycle with real database
- 21 integration tests covering:
  - `BookController`: search endpoint (Cyrillic/Latin queries, validation, limits)
  - `ParserController`: parse endpoints with mocked parser
  - `BooksRepository`: database operations (search, find by ID, rebuild)

**Test Structure:**

```go
// tests/integration/book_controller_test.go
type BookControllerSuite struct {
    suite.Suite
    db         *testhelpers.TestDatabase
    repo       *repositories.BooksRepository
    controller *controllers.BookController
    router     *gin.Engine
}
```

### Unit Tests (Vitest)

Located in `web/` directory:

- `web/pages/admin.test.ts` - Component tests
- `web/shared/api/*.test.ts` - API function tests

```bash
bun run test
```

Use AAA pattern in lowercase: arrange, act, assert

**Pattern:** HTTP mocking via `vi.mock()`:

```typescript
vi.mock('wretch', async () => {
    const actual = await vi.importActual('wretch')
    return {
        ...actual,
        default: vi.fn().mockReturnValue({
            get: vi.fn().mockResolvedValue({ fetch: vi.fn().mockResolvedValue(...) })
        })
    }
})
```

### E2E Tests (Playwright)

Located in `e2e/tests/`:

- `e2e/tests/search.test.ts` - Search flow tests
- `e2e/tests/admin.test.ts` - Admin page tests

```bash
bun run playwright
```

**Pattern:** Page Object Pattern:

```typescript
// playwright/pages/home.page.ts
class HomePage {
  constructor(private page: Page) {}

  async search(query: string) {
    await this.page.getByPlaceholder("Введите название...").fill(query);
    await this.page.getByRole("button", { name: "Найти" }).click();
  }
}
```

### Test Environment

- **Vitest**: `happy-dom` environment
- **Playwright**: Starts Vite dev server on port 5173, mocks API routes

## Docker Deployment

Multi-stage build combines frontend and backend:

```dockerfile
# Stage 1: Build frontend
FROM oven/bun:1 AS frontend-builder
# Build Vue app

# Stage 2: Build backend
FROM golang:1.25 AS backend-builder
# Build Go binary

# Stage 3: Runtime
FROM alpine:latest
# Copy assets and binary, run with tini
```

Run with Docker Compose:

```bash
docker-compose up -d  # Starts PostgreSQL
```
