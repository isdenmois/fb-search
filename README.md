# FB Search

A book search application for Flibusta library. Search and download FB2 books with a modern web interface.

## Tech Stack

- **Backend**: Go (Gin framework)
- **Frontend**: Vue + TypeScript + Vite
- **Database**: PostgreSQL
- **Build**: Bun

## Project Structure

```
.
├── app/               # INP parser for Flibusta format
├── domain/            # Domain models
├── infra/             # Database layer
├── migrations/        # SQL migrations
├── parser/            # Book parsing logic
├── public/            # Built frontend assets
├── shared/            # Shared utilities
├── views/             # HTTP controllers
├── scripts/           # Build & migration scripts
└── files/             # Book archives (FB2/ZIP)
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

## API Endpoints

- `GET /api/search?q=<query>` - Search books
- `GET /dl/:file/:path` - Download book file
- `POST /api/parse` - Rebuild database from INPX files

## License

MIT
