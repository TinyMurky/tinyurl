# Project Context

## Purpose
A URL shortener service designed for efficiency and scalability. It provides APIs to shorten URLs and redirect short URLs to their original destinations. The system leverages Bloom filters and caching to optimize performance for both write and read operations.

## Tech Stack
-   **Language:** Go (1.25.1)
-   **Databases:**
    -   SQLite (Primary storage, using `modernc.org/sqlite`)
    -   Redis (Cache and Bloom Filter, using `github.com/redis/go-redis/v9`)
-   **ID Generation:** `github.com/TinyMurky/snowflake` (Snowflake IDs)
-   **Migrations:** `github.com/golang-migrate/migrate/v4`
-   **Configuration:** `github.com/sethvargo/go-envconfig`, `github.com/joho/godotenv`
-   **Logging:** `go.uber.org/zap` (Structured logging)
-   **Containerization:** Docker, Docker Compose (Development and Debug setups)
-   **Frameworks:** Standard library `net/http` for routing.

## Project Conventions

### Code Style
-   Follows [Uber Go Style Guide](https://github.com/uber-go/guide).
-   Uses [revive](https://github.com/mgechev/revive) for linting.
-   Folder structure follows [google/exposure-notifications-server](https://github.com/google/exposure-notifications-server) pattern:
    -   `cmd/`: Entry points for applications and migrations.
    -   `internal/`: Internal logic, handlers, and service-specific code.
    -   `pkg/`: Reusable packages and library-like code.
-   Uses `net/http` standard library for routing (Go 1.22+ patterns).

### Architecture Patterns
-   **Dependency Injection:** Uses a `ServerEnv` and `Setup` pattern to wire dependencies (DB, Cache, Bloom Filter) into servers.
-   **Write Path:**
    1.  Check existence (Bloom Filter/DB).
    2.  Generate unique ID (Snowflake).
    3.  Persist to Database.
    4.  Update Bloom Filter.
    5.  Update Cache.
-   **Read Path:**
    1.  Check Bloom Filter (fail fast if definitely not present).
    2.  Check Cache (Redis).
    3.  Query Database (if cache miss).
    4.  Populate Cache (if found in DB).
    5.  Redirect (302) or 404.
-   **Middleware:** Centralized middleware stack for logging, recovery, etc.

### Testing Strategy
-   Unit tests alongside code (`_test.go`).
-   Docker Compose for integration/environment testing (`make up-debug`, `make down-debig`).

### Git Workflow
-   (Use Github workflow assumed, only main and feature branch)

## Domain Context
-   **Core Entities:** URL (Long URL, Short Key, Metadata).
-   **Endpoints:**
    -   `POST /api/v1/data/shorten`: Accepts `{longUrl: string}`, returns short URL.
    -   `GET /api/v1/shortUrl`: Redirects to original URL.
    -   `GET /shortUrl`: Alias for redirect. (planned)
    -   `GET /`: UI (planned).

## Important Constraints
-   **Performance:** High read throughput is critical.
-   **Integrity:** Bloom filters used to reduce DB load for non-existent keys.

## External Dependencies
-   **Redis:** Required for caching and Bloom filters.
-   **SQLite:** Used as the persistent data store.