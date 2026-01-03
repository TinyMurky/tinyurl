# Design: Singleflight Integration

## Context
The URL shortener is read-heavy. A cache stampede can occur when a hot key expires.

## Solution
Use `golang.org/x/sync/singleflight` wrapped in a dependency-injected package with configuration.

### Architecture
1.  **`pkg/singleflight`**: Defines `Group` interface, `Config` struct, and `NewFromEnv`.
2.  **`ServerEnv`**: Instantiates and holds a `singleflight.Group`.
3.  **Injection**: The `Handler` receives `ServerEnv` and extracts the `singleflight` group.

### Configuration
Although `x/sync/singleflight` has no inherent tuning knobs, we will add a configuration layer to support future features like:
-   **Timeout:** A default timeout for all flights (optional wrapper logic).
-   **Metrics:** Toggle for observability.
-   **Disable:** A feature flag to bypass singleflight (useful for debugging).

### Flow
1.  **Request arrives.**
2.  **Bloom/Cache:** Check (Fail fast).
3.  **Singleflight**:
    -   Handler calls `group.Do("shorturl:" + id, fn)`.
    -   **Fn**: Query DB -> Update Cache -> Return Model.
4.  **Result**: Shared result returned to all callers.

## Considerations
-   **Global vs. Local**: Using `ServerEnv` makes it a singleton. Key collision is a risk.
    -   *Mitigation:* Use strict key prefixing (e.g., `entity:id`).
