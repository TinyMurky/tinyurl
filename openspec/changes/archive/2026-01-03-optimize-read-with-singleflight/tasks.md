# Tasks

-   [ ] Add `golang.org/x/sync` to `go.mod`. <!-- id: 0 -->
-   [ ] Create `pkg/singleflight` with `Group` interface. <!-- id: 4 -->
-   [ ] Implement `pkg/singleflight/config.go` with `Config` struct and `NewFromEnv`. <!-- id: 7 -->
-   [ ] Implement `pkg/singleflight/singleflight.go` wrapping `x/sync`. <!-- id: 8 -->
-   [ ] Update `internal/serverenv` to include `pkg/singleflight.Group`. <!-- id: 5 -->
-   [ ] Update `internal/setup` to initialize `singleflight` from config. <!-- id: 6 -->
-   [ ] Update `Handler` in `handle_get_shorturl` to use `env.SingleFlight()`. <!-- id: 1 -->
-   [ ] Refactor `ServeHTTP` to use `group.Do` with `shorturl:` prefix. <!-- id: 2 -->
-   [ ] Verify with tests. <!-- id: 3 -->
