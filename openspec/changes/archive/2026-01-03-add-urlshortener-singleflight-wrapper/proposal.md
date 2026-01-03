# Add URLShortener Singleflight Wrapper

## Why
Currently, the `handle_get_shorturl` handler manually prefixes keys with "shorturl:" when using the generic `singleflight` package. This leaks implementation details (key namespacing) into the HTTP handler and risks inconsistency if other handlers need similar logic.

## What Changes
-   Create a new package `internal/urlshortener/singleflight`.
-   Implement a domain-specific wrapper around `pkg/singleflight.Group`.
-   Encapsulate the "shorturl:" key prefixing logic within this wrapper.
-   Refactor `handle_get_shorturl` to use this new wrapper.

## Impact
-   **New Package:** `internal/urlshortener/singleflight`
-   **Refactor:** `internal/urlshortener/api/v1/handle_get_shorturl/handler.go`
