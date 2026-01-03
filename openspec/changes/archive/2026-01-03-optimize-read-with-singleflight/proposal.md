# Optimize Read Path with Singleflight

## Why
Protect the database from cache stampedes (thundering herd problem) during high-concurrency access to the same short URL when it is missing from the cache.

## What Changes
-   Create `pkg/singleflight` to wrap `golang.org/x/sync/singleflight` and allow dependency injection.
-   Add `SingleFlight` to `internal/serverenv` to make it available globally.
-   Update `GET /api/v1/shortUrl` handler to use the injected `singleflight` service.
-   Ensure concurrent requests for the same ID share a single database query.

## Impact
-   **New Package:** `pkg/singleflight`
-   **Internal:** `ServerEnv` gains a new field.
-   **Handler:** `handle_get_shorturl` consumes the service.